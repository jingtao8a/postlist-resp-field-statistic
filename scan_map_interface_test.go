package main

import (
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"fmt"
	"reflect"
	"testing"
)

// 存储找到的map[string]interface{}字段信息
type MapFieldInfo struct {
	FullPath     string // 完整路径，如ParentStruct.ChildStruct.FieldName
	FieldName    string // 字段名
	ParentStruct string // 直接父结构体名
	Depth        int    // 嵌套深度
}

// 检查类型是否为map[string]interface{}
func isMapStringInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Map &&
		t.Key().Kind() == reflect.String &&
		t.Elem().Kind() == reflect.Interface
}

// 递归检查结构体字段，收集map[string]interface{}字段信息
func collectMapStringInterfaceFields(t reflect.Type, currentPath string, depth int,
	visited map[reflect.Type]bool, results *[]MapFieldInfo) {
	// 避免循环引用导致的无限递归
	if visited[t] {
		return
	}
	visited[t] = true

	// 只处理结构体类型
	if t.Kind() != reflect.Struct {
		return
	}

	structName := t.Name()
	// 处理根结构体路径
	fullStructPath := currentPath
	if currentPath == "" {
		fullStructPath = structName
	} else {
		fullStructPath = fmt.Sprintf("%s.%s", currentPath, structName)
	}

	// 遍历结构体的所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldName := field.Name
		fieldFullPath := fmt.Sprintf("%s.%s", fullStructPath, fieldName)

		// 检查当前字段是否为map[string]interface{}
		if isMapStringInterface(fieldType) {
			*results = append(*results, MapFieldInfo{
				FullPath:     fieldFullPath,
				FieldName:    fieldName,
				ParentStruct: structName,
				Depth:        depth,
			})
			continue // 找到后继续检查其他字段
		}

		// 处理指针类型
		elemType := fieldType
		if fieldType.Kind() == reflect.Ptr {
			elemType = fieldType.Elem()
		}

		// 递归处理嵌套结构体
		if elemType.Kind() == reflect.Struct {
			collectMapStringInterfaceFields(elemType, fieldFullPath, depth+1, visited, results)
		}

		// 处理切片/数组中的结构体
		if elemType.Kind() == reflect.Slice || elemType.Kind() == reflect.Array {
			sliceElemType := elemType.Elem()
			if sliceElemType.Kind() == reflect.Ptr {
				sliceElemType = sliceElemType.Elem()
			}
			if sliceElemType.Kind() == reflect.Struct {
				collectMapStringInterfaceFields(sliceElemType, fieldFullPath, depth+1, visited, results)
			}
		}
	}
}

// 对外暴露的检查函数，返回所有找到的map字段信息
func FindMapStringInterfaceFields(s interface{}) []MapFieldInfo {
	var results []MapFieldInfo
	t := reflect.TypeOf(s)

	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 确保输入是结构体
	if t.Kind() != reflect.Struct {
		return results
	}

	visited := make(map[reflect.Type]bool)
	collectMapStringInterfaceFields(t, "", 0, visited, &results)
	return results
}

// 格式化打印字段信息
func PrintMapFieldInfo(fields []MapFieldInfo) {
	if len(fields) == 0 {
		fmt.Println("未找到map[string]interface{}类型的字段")
		return
	}

	fmt.Printf("共找到 %d 个map[string]interface{}类型的字段：\n", len(fields))
	for i, field := range fields {
		fmt.Printf("\n字段 %d:\n", i+1)
		fmt.Printf("  完整路径: %s\n", field.FullPath)
		fmt.Printf("  字段名: %s\n", field.FieldName)
		fmt.Printf("  父结构体: %s\n", field.ParentStruct)
		fmt.Printf("  嵌套深度: %d\n", field.Depth)
	}
}

// 示例用法
func TestScanMapInterface(t *testing.T) {
	test := com_ss_ugc_tiktok.AwemeV1AwemePostResponse{}
	fields := FindMapStringInterfaceFields(test)
	PrintMapFieldInfo(fields)
}
