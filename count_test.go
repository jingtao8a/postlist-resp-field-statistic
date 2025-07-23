package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestCount(t *testing.T) {
	for i := 0; i < 100; i++ {
		go saveResponse()
	}
	for true {
		// do nothing
	}
}

type SaveResponseContext struct {
	mu sync.RWMutex
}

var saveResponseContext SaveResponseContext

func saveResponse() {
	saveResponseContext.mu.Lock()
	defer saveResponseContext.mu.Unlock()
	// 打开文件（追加模式，如果不存在则创建）
	countFile, err := os.OpenFile("count.txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("[saveResponse] open countFile failed: %v", err)
		return
	}
	defer func(countFile *os.File) {
		err = countFile.Close()
		if err != nil {
			fmt.Printf("close countFile failed: %v", err)
		}
	}(countFile)
	data, _ := ioutil.ReadAll(countFile)
	current, _ := strconv.Atoi(string(data))
	newValue := current + 1
	// 清空文件内容
	if err := countFile.Truncate(0); err != nil {
		fmt.Println("Error truncating file:", err)
		return
	}

	// 可选：将文件指针重置到文件开头（若后续需要写入）
	if _, err := countFile.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		return
	}
	fmt.Println(newValue)
	if _, err := countFile.Write([]byte(fmt.Sprintf("%d", newValue))); err != nil {
		fmt.Printf("[saveResponse] write countFile failed: %v", err)
	}
}
