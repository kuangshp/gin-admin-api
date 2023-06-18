package utils

import (
	"fmt"
	"reflect"
)

// CopyProperties 简单的拷贝对象的方法，dst目标的一个地址
func CopyProperties(dst, src interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(dst)
	av := reflect.ValueOf(dst)
	bt := reflect.TypeOf(src)
	bv := reflect.ValueOf(src)
	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("目标结构体必须是一个指针类型")
		return
	}
	av = reflect.ValueOf(av.Interface())
	// 要复制哪些字段
	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}
	if len(_fields) == 0 {
		fmt.Println("没有字段可拷贝")
		return
	}
	// 复制
	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)
		// a中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		} else {
			fmt.Printf("no such field or different kind, fieldName: %s\n", name)
		}
	}
	return
}
