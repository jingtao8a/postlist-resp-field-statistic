package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestJsonSerialization(t *testing.T) {
	type NestedStruct struct {
		A *string
		B *string
	}
	type Data struct {
		Name   *string
		Age    *int
		Nested *NestedStruct
	}
	str := "yuxintao"
	age := 12
	data := &Data{
		Name:   &str,
		Age:    &age,
		Nested: &NestedStruct{},
	}
	// 序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
	// 输出：{"active":true,"age":25,"hobbies":["swimming","hiking"],"name":"Bob"}
}

func TestJsonReadToFile(t *testing.T) {
	data := map[string]interface{}{
		"name":    "Bob",
		"age":     25,
		"hobbies": []string{"swimming", "hiking"},
		"active":  true,
	}
	// 序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var strs []string
	for i := 0; i < 10; i++ {
		strs = append(strs, string(jsonData))
	}
	file, err := os.Create("test.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()
	for i := 0; i < len(strs); i++ {
		file.WriteString(strs[i] + "\n")
	}
}

func TestJsonReadFromFile(t *testing.T) {
	// 打开文件
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	var strs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}
	fmt.Println(len(strs))
	for i := 0; i < len(strs); i++ {
		fmt.Println(strs[i])
	}
}
