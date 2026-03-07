package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CampingTimeSlot 预约时段
type CampingTimeSlot struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;comment:时段名称;size:50;"`
	StartTime string `json:"startTime" form:"startTime" gorm:"column:start_time;comment:开始时间如08:00;size:10;"`
	EndTime   string `json:"endTime" form:"endTime" gorm:"column:end_time;comment:结束时间如12:00;size:10;"`
	Sort      int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0;"`
}

// TableName 表名
func (CampingTimeSlot) TableName() string {
	return "gva_camping_time_slots"
}
