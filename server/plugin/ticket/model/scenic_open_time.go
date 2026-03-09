package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ScenicOpenTime 景区开放时间（按星期）
type ScenicOpenTime struct {
	global.GVA_MODEL
	ScenicID  uint     `json:"scenicId" form:"scenicId" gorm:"column:scenic_id;comment:景区ID;"`
	WeekDay   int      `json:"weekDay" form:"weekDay" gorm:"column:week_day;comment:星期1-7;"`
	OpenTime  TimeOnly `json:"openTime" form:"openTime" gorm:"column:open_time;comment:开放时间;"`
	CloseTime TimeOnly `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;"`
}

// TableName 表名
func (ScenicOpenTime) TableName() string {
	return "scenic_open_time"
}
