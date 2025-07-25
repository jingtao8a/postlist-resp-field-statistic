package main

import (
	"fmt"
	"log"
	"postlist-resp-field-statistic/common"
	"postlist-resp-field-statistic/config"
	"reflect"
	"testing"
)

type Inner struct {
	Data map[string]interface{}
}

type Middle struct {
	Inners []*Inner
}

type Outer struct {
	Name    string
	Middle  *Middle
	Meta    map[string]interface{}
	History []map[string]interface{}
	Extras  map[string]*Inner
}

func TestScan(t *testing.T) {
	responses, err := common.ReadMessagesFromFile(config.RESPONSE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	if responses == nil {
		log.Fatal("responses is nil")
	}
	var totalResponseCount = len(responses)
	fmt.Printf("get %d responses int total\n", totalResponseCount)
	fmt.Printf("=================================\n")
	for _, response := range responses {
		findMapStringInterfaceFieldsWithValue(reflect.ValueOf(response), "AwemeV1AwemePostResponse")
	}
}

// 递归查找并打印字段是否赋值
func findMapStringInterfaceFieldsWithValue(v reflect.Value, path string) {
	if !v.IsValid() {
		return
	}

	// 处理指针，拿指针指向的值
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// nil 指针，不能递归值，只能递归类型判断
			break
		}
		v = v.Elem()
	}

	t := v.Type()
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			fieldVal := v.Field(i)
			fieldPath := path + "." + field.Name
			findMapStringInterfaceFieldsWithValue(fieldVal, fieldPath)
		}

	case reflect.Slice, reflect.Array:
		// 遍历每个元素（如果有的话）
		for i := 0; i < v.Len(); i++ {
			elemVal := v.Index(i)
			findMapStringInterfaceFieldsWithValue(elemVal, fmt.Sprintf("%s[%d]", path, i))
		}

	case reflect.Map:
		// 判断是否是 map[string]interface{}
		if t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.Interface {
			// map[string]interface{} 赋值判断
			isSet := !v.IsNil() && v.Len() > 0
			fmt.Printf("字段：%s 类型：map[string]interface{} 是否赋值：%v\n", path, isSet)
			return
		}
		// 不是 map[string]interface{}，递归遍历 map 的所有 value
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			findMapStringInterfaceFieldsWithValue(val, fmt.Sprintf("%s[%v]", path, key))
		}
	}
}
