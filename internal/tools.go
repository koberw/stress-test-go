package tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetDirPath() string {
	_, filename, _, ok := runtime.Caller(0) // 0 表示当前函数所在文件
	if !ok {
		panic("无法获取当前文件信息")
	}

	// 获取可执行文件所在的目录
	exeDir := filepath.Dir(filename)
	fmt.Println("当前目录:", exeDir)
	return exeDir
}

func GetRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("项目根目录: %s\n", dir)
	return dir
}

func GetConfigsPath() string {
	rootPath := GetRootPath()
	path := rootPath + "/configs"
	fmt.Printf("configs目录: %s\n", path)
	return path
}

func GetReportsPath() string {
	rootPath := GetRootPath()
	path := rootPath + "/reports"
	fmt.Printf("reports目录: %s\n", path)
	return path
}
