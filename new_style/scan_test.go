package main

import (
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"fmt"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	findAllMapStringInterfaceFields(reflect.TypeOf(com_ss_ugc_tiktok.AwemeV1AwemePostResponse{}))
}

func findAllMapStringInterfaceFields(t reflect.Type) {
	visited := map[reflect.Type]bool{}
	recursiveFind(t, "AwemeV1AwemePostResponse", visited)
}

func recursiveFind(t reflect.Type, path string, visited map[reflect.Type]bool) {
	if t == nil {
		return
	}

	// 先把指针解引用
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t == nil {
		return
	}

	switch t.Kind() {
	case reflect.Struct:
		// 防止递归死循环
		if visited[t] {
			return
		}
		visited[t] = true

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// 只遍历导出字段
			if !field.IsExported() {
				continue
			}
			fieldPath := path + "." + field.Name
			recursiveFind(field.Type, fieldPath, visited)
		}

	case reflect.Slice, reflect.Array:
		elemType := t.Elem()
		recursiveFind(elemType, path+"[]", visited)

	case reflect.Map:
		keyType := t.Key()
		valType := t.Elem()

		// 检查是否是 map[string]interface{}
		if keyType.Kind() == reflect.String && valType.Kind() == reflect.Interface {
			fmt.Printf("找到字段：%s，类型：map[string]interface{}\n", path)
		}
		// 继续递归 map 的值类型
		recursiveFind(valType, path+"[key]", visited)
	}
}
