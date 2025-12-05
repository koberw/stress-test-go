package output

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var OutputDirPath string

var (
	recordsFile *os.File   // 全局文件句柄
	statsFile   *os.File   // 全局文件句柄
	once        sync.Once  // 保证只初始化一次
	mu          sync.Mutex // 保证并发写入安全
)

func Init(testId string) {
	OutputDirPath = getDirPath()
	testTaskPath := OutputDirPath + "/" + testId
	recordsPath := testTaskPath + "/records"
	fmt.Println("outpur records -> " + recordsPath)
	statsPath := testTaskPath + "/stats"
	fmt.Println("outpur stats -> " + statsPath)

	once.Do(func() {
		//创建testid目录
		err := os.MkdirAll(testTaskPath, 0755)
		if err != nil {
			fmt.Println("创建目录失败:", err)
			return
		}

		// 以追加模式打开文件（如果不存在则创建），并且可写
		recordsFile, err = os.OpenFile(recordsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		statsFile, err = os.OpenFile(statsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
}

func getDirPath() string {
	_, filename, _, ok := runtime.Caller(0) // 0 表示当前函数所在文件
	if !ok {
		panic("无法获取当前文件信息")
	}

	// 获取可执行文件所在的目录
	exeDir := filepath.Dir(filename)
	fmt.Println("output根目录:", exeDir)
	return exeDir
}

func WriteLineToRecords(line string) error {
	mu.Lock()
	defer mu.Unlock()

	if recordsFile == nil {
		fmt.Println("recordsFile == nil")
		return fmt.Errorf("文件未初始化，请先调用 Init(testId)")
	}

	// 写入内容 + 换行
	_, err := recordsFile.WriteString(line + "\n")
	return err
}

func WriteLineToStats(line string) error {
	if statsFile == nil {
		return fmt.Errorf("文件未初始化，请先调用 Init(testId)")
	}

	// 写入内容 + 换行
	_, err := statsFile.WriteString(line + "\n")
	return err
}

// CloseFile 关闭文件（应该在程序退出前调用）
func CloseFile() {
	mu.Lock()
	defer mu.Unlock()

	if recordsFile != nil {
		recordsFile.Close()
	}

	if statsFile != nil {
		statsFile.Close()
	}
}
