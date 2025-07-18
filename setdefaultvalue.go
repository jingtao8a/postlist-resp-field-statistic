package main

import (
	"reflect"
)

// SetDefaultValues 递归设置结构体中所有 nil 指针字段的默认值
func SetDefaultValues(v interface{}) {
	val := reflect.ValueOf(v)

	// 必须传入指针，否则无法修改原始值
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}

	// 解引用指针，获取实际结构体
	elem := val.Elem()
	if !elem.CanSet() {
		return
	}

	// 处理结构体类型
	if elem.Kind() == reflect.Struct {
		setStructDefaults(elem)
	}
}

// setStructDefaults 处理结构体的每个字段
func setStructDefaults(structVal reflect.Value) {
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		// 忽略未导出的私有字段
		if fieldType.PkgPath != "" {
			continue
		}

		// 处理指针类型
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				setDefaultForNilPtr(field, fieldType.Type.Elem())
			} else if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				// 递归处理嵌套结构体指针
				setStructDefaults(field.Elem())
			}
		}

		// 处理嵌套结构体
		if field.Kind() == reflect.Struct {
			setStructDefaults(field)
		}

		// 处理切片和映射
		if field.Kind() == reflect.Slice || field.Kind() == reflect.Map {
			if field.IsNil() {
				setDefaultForCollection(field)
			}
		}
	}
}

// setDefaultForNilPtr 为 nil 指针设置默认值
func setDefaultForNilPtr(field reflect.Value, elemType reflect.Type) {
	// 创建新的实例
	newVal := reflect.New(elemType)

	// 根据元素类型设置默认值
	if elemType.Kind() == reflect.Struct {
		// 递归处理嵌套结构体
		setStructDefaults(newVal.Elem())
	}

	// 设置指针值
	field.Set(newVal)
}

// setDefaultForCollection 为 nil 切片/映射设置默认值
func setDefaultForCollection(field reflect.Value) {
	switch field.Kind() {
	case reflect.Slice:
		sliceType := field.Type()
		// 创建空切片
		newSlice := reflect.MakeSlice(sliceType, 0, 0)
		field.Set(newSlice)
	case reflect.Map:
		mapType := field.Type()
		// 创建空映射
		newMap := reflect.MakeMap(mapType)
		field.Set(newMap)
	}
}

// EnsureNestedMessages 确保所有嵌套消息指针非 nil
func EnsureNestedMessages(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}

	// 解引用指针
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	// 递归遍历结构体字段
	ensureStructFields(elem)
}

func ensureStructFields(structVal reflect.Value) {
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		// 忽略未导出字段
		if fieldType.PkgPath != "" {
			continue
		}

		// 处理嵌套消息指针（结构体指针）
		if field.Kind() == reflect.Ptr && !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			ensureStructFields(field.Elem()) // 递归处理嵌套结构体
		}

		// 初始化 nil 的嵌套消息指针
		if field.Kind() == reflect.Ptr && field.IsNil() && field.Type().Elem().Kind() == reflect.Struct {
			newVal := reflect.New(field.Type().Elem())
			field.Set(newVal)
			ensureStructFields(newVal.Elem()) // 递归初始化新创建的结构体
		}
	}
}
