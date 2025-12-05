package output

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RunStats(testId string, startTime time.Time, endTime time.Time) {
	successNum := 0
	costSum := 0
	var costList []int

	// 使用 bufio.Scanner 逐行读取
	file, err := os.Open(recordsFile.Name())
	defer file.Close()

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行内容（不含换行符）
		columns := strings.Split(line, "\t")
		if len(columns) != 5 {
			continue
		}
		success, _ := strconv.ParseBool(columns[2])
		if success {
			successNum++
		}

		cost, _ := strconv.Atoi(columns[3])
		costList = append(costList, cost) // 将该行加入切片
		costSum += cost
	}
	sort.Ints(costList)

	//output
	WriteLineToStats("-------------- 压测统计报告 ---------------")
	WriteLineToStats("[test id]" + testId)
	WriteLineToStats("[test start]" + startTime.Format("2006-01-02 15:04:05"))
	WriteLineToStats("[test end]" + endTime.Format("2006-01-02 15:04:05"))
	testDuration := endTime.Sub(startTime).Seconds()
	WriteLineToStats("[test duration(s)]" + fmt.Sprintf("%v", testDuration))
	totalRequestNum := len(costList)
	WriteLineToStats("[total requests]" + strconv.Itoa(totalRequestNum))
	WriteLineToStats("[success num]" + strconv.Itoa(successNum))
	percent := (float64(successNum) / float64(totalRequestNum)) * 100
	WriteLineToStats("[success rate(%)]" + fmt.Sprintf("%0.2f", percent) + "%")
	WriteLineToStats("[cost min(ms)]" + strconv.Itoa(costList[0]))
	WriteLineToStats("[cost max(ms)]" + strconv.Itoa(costList[len(costList)-1]))
	avgCost := costSum / len(costList)
	WriteLineToStats("[cost avg(ms)]" + strconv.Itoa(avgCost))
	index99 := int(math.Floor(float64(len(costList))*0.99)) - 1
	WriteLineToStats("[cost P99(ms)]" + strconv.Itoa(costList[index99]))
	fmt.Println(">>>>>>>>>>> 请查看压测报告:" + statsFile.Name())
}
