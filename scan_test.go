package main

import (
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"fmt"
	"reflect"
	"testing"
)

// 扫描结果
type ScanResult struct {
	FieldPath string       // 字段路径（如 "User.Info.Metadata"）
	FieldType reflect.Type // 字段类型
	Value     interface{}  // 字段值
}

// ScanInterfaceFields 递归扫描结构体，找出所有 interface{} 类型的字段
func ScanInterfaceFields(obj interface{}) []ScanResult {
	var results []ScanResult
	scanValue(reflect.ValueOf(obj), "", &results)
	return results
}

// 递归扫描值
func scanValue(v reflect.Value, path string, results *[]ScanResult) {
	// 处理指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		scanValue(v.Elem(), path, results)
		return
	}

	// 处理接口
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		scanValue(v.Elem(), path, results)
		return
	}

	// 处理结构体
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)

			// 构建完整路径
			fieldPath := fieldType.Name
			if path != "" {
				fieldPath = path + "." + fieldPath
			}

			// 如果字段是匿名字段（嵌入结构体），不添加字段名
			if fieldType.Anonymous {
				scanValue(field, path, results)
			} else {
				// 检查字段类型是否为 interface{}
				if fieldType.Type.Kind() == reflect.Interface &&
					fieldType.Type.String() == "interface {}" {
					*results = append(*results, ScanResult{
						FieldPath: fieldPath,
						FieldType: fieldType.Type,
						Value:     field.Interface(),
					})
				}

				// 递归扫描嵌套字段
				scanValue(field, fieldPath, results)
			}
		}
		return
	}

	// 处理切片和数组
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			// 构建路径（如 "Items[0]"）
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			scanValue(v.Index(i), itemPath, results)
		}
		return
	}

	// 处理映射
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			// 构建路径（如 "Data[key]"）
			keyStr := fmt.Sprintf("%v", key.Interface())
			itemPath := fmt.Sprintf("%s[%s]", path, keyStr)
			scanValue(v.MapIndex(key), itemPath, results)
		}
		return
	}
}

// 示例使用
func TestScanInterfaceField(t *testing.T) {
	// 扫描 interface{} 字段
	results := ScanInterfaceFields(&com_ss_ugc_tiktok.AwemeV1AwemePostResponse{})

	// 输出结果
	if len(results) == 0 {
		fmt.Println("未发现 interface{} 类型的字段")
	} else {
		fmt.Println("发现以下 interface{} 类型的字段：")
		for _, r := range results {
			fmt.Printf("- 路径: %s\n", r.FieldPath)
			fmt.Printf("  类型: %s\n", r.FieldType)
			fmt.Printf("  值: %v\n", r.Value)
			fmt.Println()
		}
	}
}
