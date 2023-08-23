package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Decimal = decimal.Decimal

// LocalTime 自定义数据类型1开始
type LocalTime struct {
	time.Time
}

// MarshalJSON 序列化为JSON
func (t LocalTime) MarshalJSON() ([]byte, error) {
	//TODO 返回时间戳 格式化秒
	seconds := t.UnixNano() / 1e6 // 毫秒时间
	if seconds > 0 {
		return []byte(fmt.Sprintf("%d", seconds)), nil
	} else {
		return []byte("\"\""), nil
	}
	// TODO 如果要返回字符串格式化数据
	//if t.IsZero() {
	//	return []byte("\"\""), nil
	//}
	//stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	//return []byte(stamp), nil
}

// UnmarshalJSON 反序列化为JSON
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var err error
	t.Time, err = time.Parse("2006-01-02 15:04:05", string(data)[1:20])
	return err
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
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

// String 重写String方法
func (t *LocalTime) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// SetRaw 读取数据库值
func (t *LocalTime) SetRaw(value interface{}) error {
	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	}
	return nil
}

// RawValue 写入数据库
func (t *LocalTime) RawValue() interface{} {
	str := t.Format("2006-01-02 15:04:05")
	if str == "0001-01-01 00:00:00" {
		return nil
	}
	return str
}

// LocalTime 自定义数据类型1结束

// JSON 自定义JSON数据类型
type JSON []byte

func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		errors.New("invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}
func (m JSON) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}
func (m *JSON) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("null point exception")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}
func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}
