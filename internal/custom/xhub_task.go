package custom

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	tools "stress_test-go/internal"
	"stress_test-go/internal/output"
)

type XhubTask struct {
	queryList []string
}

// 构造函数，初始化queryList
func NewXhubTask() *XhubTask {
	return &XhubTask{
		queryList: loadQueries(),
	}
}

func loadQueries() []string {
	queryPath := tools.GetConfigsPath() + "/" + "query.txt"
	log.Printf("queryPath: %s", queryPath)

	// 打开文件
	file, err := os.Open(queryPath) // 替换为你的文件路径
	if err != nil {
		log.Fatal("无法打开文件:", err)
	}
	defer file.Close() // 确保文件最终关闭

	var queryList []string
	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()              // 获取当前行内容（不含换行符）
		queryList = append(queryList, line) // 将该行加入切片
	}

	// 检查扫描过程中是否有错误（如文件损坏）
	if err := scanner.Err(); err != nil {
		log.Fatal("读取文件出错:", err)
	}

	// lines 就是你要的 []string，可以后续用作其他处理
	fmt.Println("\nquery总共读取了", len(queryList), "行")

	return queryList
}

// 执行单个 HTTP 请求
func (xhub *XhubTask) DoRequest() {
	requestId := uuid.New().String()
	url := "http://xhub.xsearch.woa.com/tianji/SearchPassage"

	//随机选择一个query
	randomIndex := rand.Intn(len(xhub.queryList))
	query := xhub.queryList[randomIndex]
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
