package input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var QueryList []string

func init() {
	loadQueries()
}

func getDirPath() string {
	_, filename, _, ok := runtime.Caller(0) // 0 表示当前函数所在文件
	if !ok {
		panic("无法获取当前文件信息")
	}

	// 获取可执行文件所在的目录
	exeDir := filepath.Dir(filename)
	fmt.Println("input根目录:", exeDir)
	return exeDir
}

func loadQueries() {
	queryPath := getDirPath() + "/" + "query.txt"
	log.Printf("queryPath: %s", queryPath)

	// 打开文件
	file, err := os.Open(queryPath) // 替换为你的文件路径
	if err != nil {
		log.Fatal("无法打开文件:", err)
	}
	defer file.Close() // 确保文件最终关闭

	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()              // 获取当前行内容（不含换行符）
		QueryList = append(QueryList, line) // 将该行加入切片
	}

	// 检查扫描过程中是否有错误（如文件损坏）
	if err := scanner.Err(); err != nil {
		log.Fatal("读取文件出错:", err)
	}

	// lines 就是你要的 []string，可以后续用作其他处理
	fmt.Println("\n总共读取了", len(QueryList), "行")
}
