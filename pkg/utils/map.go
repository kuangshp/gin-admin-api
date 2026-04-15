package utils

import (
	"encoding/json"
	"fmt"
	"sort"
)

func MapToJson(result interface{}) string {
	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}

// MapKeySort 将map的key进行ASCII排序
func MapKeySort[T int8 | int16 | int32 | int64 | int | float32 | float64 | string | bool | interface{}](m map[string]T, options ...string) string {
	// 将map中全部的key到切片中
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// 对切片进行排序
	sort.Strings(keys)
	// 循环拼接返回
	result := ""
	for _, item := range keys {
		if len(options) > 0 {
			if result != "" {
				result += fmt.Sprintf("%s%v", options[0], fmt.Sprintf("%s=%v", item, m[item]))
			} else {
				result += fmt.Sprintf("%s=%v", item, m[item])
			}
		} else {
			if result != "" {
				result += fmt.Sprintf("&%v", fmt.Sprintf("%s=%v", item, m[item]))
			} else {
				result += fmt.Sprintf("%s=%v", item, m[item])
			}
		}
	}
	return result
}
