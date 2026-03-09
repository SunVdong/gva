package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// TimeOnly 仅时间（HH:mm 或 HH:mm:ss），用于 MySQL TIME 列
type TimeOnly string

// Scan 实现 sql.Scanner，从 DB 读取 TIME
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = TimeOnly(v.Format("15:04:05"))
		return nil
	case []byte:
		*t = TimeOnly(string(v))
		return nil
	case string:
		*t = TimeOnly(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", value)
	}
}

// Value 实现 driver.Valuer。返回 "HH:mm:ss" 字符串，仅适用于 MySQL TIME 列
func (t TimeOnly) Value() (driver.Value, error) {
	s := string(t)
	if s == "" {
		return nil, nil
	}
	if len(s) == 5 && s[2] == ':' {
		return s + ":00", nil
	}
	return s, nil
}

// MarshalJSON 输出 "HH:mm"
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	s := string(t)
	if s == "" {
		return json.Marshal("")
	}
	if len(s) > 5 {
		s = s[:5]
	}
	return json.Marshal(s)
}

// UnmarshalJSON 支持 "HH:mm" 或 "HH:mm:ss"
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) == 5 && s[2] == ':' {
		*t = TimeOnly(s + ":00")
		return nil
	}
	*t = TimeOnly(s)
	return nil
}
