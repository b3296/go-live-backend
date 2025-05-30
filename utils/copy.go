package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// CopyOptions 复制选项
type CopyOptions struct {
	ExcludeFields []string               // 不复制的字段名列表
	ExtraFields   map[string]interface{} // 额外赋值字段，字段名=>值
}

// Contains 判断字符串切片是否包含某字符串
func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// CopyStructExt 支持排除字段和额外字段赋值的结构体复制
func CopyStructExt(src interface{}, dst interface{}, opts *CopyOptions) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return errors.New("参数必须是指针类型")
	}

	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()

	if !srcVal.IsValid() || !dstVal.IsValid() {
		return errors.New("无效的结构体")
	}

	srcType := srcVal.Type()

	// 先复制同名同类型字段（排除掉 ExcludeFields）
	for i := 0; i < srcVal.NumField(); i++ {
		field := srcType.Field(i)
		name := field.Name

		if opts != nil && Contains(opts.ExcludeFields, name) {
			continue
		}

		srcField := srcVal.FieldByName(name)
		dstField := dstVal.FieldByName(name)

		if dstField.IsValid() && dstField.CanSet() {
			if trySetValue(dstField, srcField) {
				continue
			}
		}
	}

	// 再赋值额外字段
	if opts != nil && opts.ExtraFields != nil {
		for key, val := range opts.ExtraFields {
			dstField := dstVal.FieldByName(key)
			if dstField.IsValid() && dstField.CanSet() {
				setInterfaceValue(dstField, val)
			}
		}
	}

	return nil
}

// 尝试给dstField赋srcField值（支持基本类型转换）
func trySetValue(dstField, srcField reflect.Value) bool {
	if srcField.Type() == dstField.Type() {
		dstField.Set(srcField)
		return true
	}

	// 简单处理 int 和 int64 转换
	if dstField.Kind() == reflect.Int && srcField.Kind() == reflect.Int64 {
		dstField.SetInt(srcField.Int())
		return true
	}
	if dstField.Kind() == reflect.Int64 && srcField.Kind() == reflect.Int {
		dstField.SetInt(srcField.Int())
		return true
	}

	// 还可以加更多类型转换

	return false
}

// 通过 interface{} 给字段赋值
func setInterfaceValue(dstField reflect.Value, val interface{}) {
	valVal := reflect.ValueOf(val)

	// 类型匹配直接赋值
	if valVal.Type().AssignableTo(dstField.Type()) {
		dstField.Set(valVal)
		return
	}

	// 类型不匹配，尝试简单转换
	switch dstField.Kind() {
	case reflect.String:
		dstField.SetString(fmt.Sprintf("%v", val))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := val.(type) {
		case int:
			dstField.SetInt(int64(v))
		case int64:
			dstField.SetInt(v)
		case float64:
			dstField.SetInt(int64(v))
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				dstField.SetInt(i)
			}
		}
	case reflect.Float32, reflect.Float64:
		switch v := val.(type) {
		case float32:
			dstField.SetFloat(float64(v))
		case float64:
			dstField.SetFloat(v)
		case int:
			dstField.SetFloat(float64(v))
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				dstField.SetFloat(f)
			}
		}
	case reflect.Bool:
		switch v := val.(type) {
		case bool:
			dstField.SetBool(v)
		case string:
			if v == "true" {
				dstField.SetBool(true)
			} else if v == "false" {
				dstField.SetBool(false)
			}
		}
		// 可以继续扩展支持更多类型
	}
}
