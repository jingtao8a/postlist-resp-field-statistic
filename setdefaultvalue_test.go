package main

import (
	"fmt"
	"testing"
)

type Nested struct {
	Value *int `json:"value"`
}

type Data struct {
	Name       *string `json:"name"`
	Age        *int    `json:"age"`
	IsValid    *bool   `json:"is_valid"`
	Nested     *Nested `json:"nested"` // 嵌套结构体指针
	Strs       []string
	NestedList []*Nested
	Map1       map[int]string
}

func TestSetDefaultValue(t *testing.T) {
	// 初始化一个包含 nil 字段的结构体
	data := &Data{}
	// 填充默认值
	SetDefaultValues(data)

	// 验证结果
	fmt.Println(data.Map1)
}
