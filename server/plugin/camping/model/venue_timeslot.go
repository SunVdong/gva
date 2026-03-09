package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VenueTimeslot 场地时间段（如 9:00-13:00、13:00-18:00）
// start_time/end_time 必须为 MySQL 的 TIME 类型；若当前为 DATETIME 会错误存储为 2008-00-00，请执行：
// ALTER TABLE venue_timeslots MODIFY start_time TIME, MODIFY end_time TIME;
type VenueTimeslot struct {
	global.GVA_MODEL
	VenueID   uint     `json:"venueId" form:"venueId" gorm:"column:venue_id;comment:场地ID;"`
	StartTime TimeOnly `json:"startTime" form:"startTime" gorm:"column:start_time;comment:开始时间;type:TIME;"`
	EndTime   TimeOnly `json:"endTime" form:"endTime" gorm:"column:end_time;comment:结束时间;type:TIME;"`
	Capacity  int      `json:"capacity" form:"capacity" gorm:"column:capacity;comment:可预约数量;default:1;"`
}

// TableName 表名
func (VenueTimeslot) TableName() string {
	return "venue_timeslots"
}
