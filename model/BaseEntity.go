package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Decimal = decimal.Decimal

// LocalTime 自定义数据类型1开始
type LocalTime struct {
	time.Time
}

// MarshalJSON 返回给前端的时候
func (t LocalTime) MarshalJSON() ([]byte, error) {
	//格式化秒
	seconds := t.UnixNano() / 1e6 // 毫秒时间
	if seconds > 0 {
		return []byte(fmt.Sprintf("%d", seconds)), nil
	} else {
		return []byte("-1"), nil
	}
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	fmt.Println(t.Time, "时间22")
	return t.Time, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// LocalTime 自定义数据类型1结束

// LocalList 自定义数据类型2开始
type LocalList []string

// Value 存储数据的时候将数组转换为字符串
func (l LocalList) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan 读取数据的时候将json字符串转换为字符串
func (l *LocalList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &l)
}
