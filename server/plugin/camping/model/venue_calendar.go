package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VenueCalendar 场地预约日历（控制某一天是否可预约）
type VenueCalendar struct {
	global.GVA_MODEL
	VenueID uint      `json:"venueId" form:"venueId" gorm:"column:venue_id;comment:场地ID;"`
	Date    time.Time `json:"date" form:"date" gorm:"column:date;comment:日期;type:date;"`
	Status  int       `json:"status" form:"status" gorm:"column:status;comment:状态1可预约0关闭;default:1;"`
}

// TableName 表名
func (VenueCalendar) TableName() string {
	return "venue_calendar"
}
