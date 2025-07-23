package main

import (
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"fmt"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	processed := make(map[reflect.Type]bool)
	scanStruct(&com_ss_ugc_tiktok.AwemeV1AwemePostResponse{}, "", processed)
}

// 扫描结构体，添加processed参数跟踪已处理类型
func scanStruct(s interface{}, parentPath string, processed map[reflect.Type]bool) {
	t := reflect.TypeOf(s)

	// 解析指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		if parentPath != "" {
			parentPath += "*"
		}
	}

	// 检查是否已处理过该类型（避免循环引用）
	if processed[t] {
		return
	}
	processed[t] = true // 标记为已处理

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			currentPath := getCurrentPath(parentPath, t.Name(), field.Name)

			if isMapStringInterfaceV2(fieldType) {
				printFieldInfo(field, fieldType, currentPath)
			} else {
				processNestedType(fieldType, currentPath, processed)
			}
		}
	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		processNestedType(t, parentPath, processed)
	}
}

// 处理嵌套类型，传递processed参数
func processNestedType(fieldType reflect.Type, currentPath string, processed map[reflect.Type]bool) {
	if fieldType.Kind() == reflect.Ptr {
		elemType := fieldType.Elem()
		if elemType.Kind() == reflect.Struct {
			scanStruct(reflect.New(elemType).Interface(), currentPath+"*", processed)
		} else if elemType.Kind() == reflect.Ptr {
			processNestedType(elemType, currentPath+"*", processed)
		}
	} else if fieldType.Kind() == reflect.Struct {
		scanStruct(reflect.New(fieldType).Interface(), currentPath, processed)
	} else if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
		elemType := fieldType.Elem()
		slicePath := currentPath + "[]"

		if elemType.Kind() == reflect.Struct {
			scanStruct(reflect.New(elemType).Interface(), slicePath, processed)
		} else if elemType.Kind() == reflect.Ptr {
			ptrElemType := elemType.Elem()
			if ptrElemType.Kind() == reflect.Struct {
				scanStruct(reflect.New(ptrElemType).Interface(), slicePath+"*", processed)
			}
		}
	}
}

// 工具函数（与之前相同）
func getCurrentPath(parentPath, structName, fieldName string) string {
	if parentPath == "" {
		return structName + "." + fieldName
	}
	return parentPath + "." + fieldName
}

func isMapStringInterfaceV2(t reflect.Type) bool {
	if t.Kind() != reflect.Map {
		return false
	}
	return t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.Interface
}

func printFieldInfo(field reflect.StructField, fieldType reflect.Type, currentPath string) {
	fmt.Printf("发现map[string]interface{}类型字段:\n")
	fmt.Printf("  字段名: %s\n", field.Name)
	fmt.Printf("  类型: %s\n", fieldType.String())
	fmt.Printf("  完整路径: %s\n", currentPath)
	fmt.Println("---")
}
