package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"stress_test-go/internal/custom"
	"stress_test-go/internal/output"
)

var (
	//测试id
	testId string
	//并发线程数, 默认200
	threadCount = 200
	//qps
	qps int
	//压测时间
	durationSec   int
	testStartTime time.Time
	testEndTime   time.Time

	customTask custom.CustomTaskRunner
)

// 压测主函数
func runTask(concurrency int, qps int, durationSec int) {
	//创建一个指定压测时间为超时时间的context
	duration := time.Duration(durationSec) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	//创建一个定时器，定时间隔为delayBetweenRequests
	delayBetweenRequests := time.Second / time.Duration(qps)
	ticker := time.NewTicker(delayBetweenRequests)
	defer ticker.Stop()

	sem := make(chan struct{}, concurrency) // 控制并发数
	var wg sync.WaitGroup

	// 用于优雅退出
	done := make(chan struct{})
	defer close(done)

	testStartTime = time.Now()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case sem <- struct{}{}: // 获取信号量
					wg.Add(1)
					go func() {
						defer wg.Done()
						defer func() { <-sem }() // 释放信号量
						//执行请求
						customTask.DoRequest()
					}()
				default:
					// 如果并发已满，跳过此次 tick（即控制 QPS 和并发双重限制）
				}
			case <-done:
				return
			}
		}
	}()

	// 等待压测时间结束
	<-ctx.Done()
	wg.Wait()

	testEndTime = time.Now()
	fmt.Println(">>>>>>>>>>>>>>>>>>>>> test complete >>>>>>>>>>>>>>>>>>>")

	time.Sleep(5 * time.Second)

	//统计
	output.RunStats(testId, testStartTime, testEndTime)

	output.CloseFile()
}

func main() {
	//定义压测配置
	qps = 5
	durationSec = 5

	fmt.Printf(">>>>>>> 开始压测:\n并发数: %d\nQPS: %d\n持续时间: %v\n >>>>>\n", threadCount, qps, durationSec)

	timestamp := time.Now().Unix()
	testId = strconv.FormatInt(timestamp, 10)
	output.Init(testId)
	fmt.Printf("测试ID：" + testId)

	//自定义任务
	customTask = custom.NewXhubTask()

	runTask(threadCount, qps, durationSec)
}
