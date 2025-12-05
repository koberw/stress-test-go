package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"stress_test_go/input"
	"stress_test_go/output"
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
)

// 执行单个 HTTP 请求
func doRequest() {
	requestId := uuid.New().String()
	url := "http://xhub.xsearch.woa.com/tianji/SearchPassage"

	//随机选择一个query
	randomIndex := rand.Intn(len(input.QueryList))
	query := input.QueryList[randomIndex]
	//fmt.Printf("use query: %v\t%s\n", randomIndex, query)

	requestBody := "{\"requestHeader\":{\"requestId\":\"" + requestId + "\",\"sessionId\":\"" + requestId +
		"\",\"guid\":\"10004\",\"scene\":{\"id\":\"10004\"}},\"queryGroup\":{\"originQuery\":\"" + query +
		"\",\"allQueries\":[{\"query\":\"" + query + "\"}]},\"requestType\":\"FROM_COMMON_SEARCH\",\"extends\":{}}"

	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		log.Fatal("创建请求失败:", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	start := time.Now()

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("请求发送失败:", err)
	}
	defer resp.Body.Close()

	var success bool = (resp.StatusCode == 200)
	log.Printf("%s -> %v\n", query, success)
	// bodyBytes, err := io.ReadAll(resp.Body)
	// log.Println(string(bodyBytes))

	cost := time.Since(start).Milliseconds()

	record := fmt.Sprintf("%s\t%s\t%v\t%v\t%s", start.Format("2006-01-02 15:04:05.000"), requestId, success, cost, query)
	output.WriteLineToRecords(record)

	io.Copy(io.Discard, resp.Body) // 丢弃响应体
	resp.Body.Close()
}

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
						doRequest()
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
	qps = 30
	durationSec = 10

	fmt.Printf(">>>>>>> 开始压测:\n并发数: %d\nQPS: %d\n持续时间: %v\n", threadCount, qps, durationSec)

	timestamp := time.Now().Unix()
	testId = strconv.FormatInt(timestamp, 10)
	output.Init(testId)

	runTask(threadCount, qps, durationSec)
}
