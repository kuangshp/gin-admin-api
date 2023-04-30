package utils

import (
	"encoding/json"
)

func MapToJson(result interface{}) string {
	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}
