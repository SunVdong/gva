package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VenueOpenTime 场地开放时间（按星期）
type VenueOpenTime struct {
	global.GVA_MODEL
	VenueID   uint   `json:"venueId" form:"venueId" gorm:"column:venue_id;comment:场地ID;"`
	WeekDay   int    `json:"weekDay" form:"weekDay" gorm:"column:week_day;comment:星期1-7;"`
	OpenTime  string `json:"openTime" form:"openTime" gorm:"column:open_time;comment:开放时间;type:time;"`
	CloseTime string `json:"closeTime" form:"closeTime" gorm:"column:close_time;comment:关闭时间;type:time;"`
}

// TableName 表名
func (VenueOpenTime) TableName() string {
	return "venue_open_time"
}
