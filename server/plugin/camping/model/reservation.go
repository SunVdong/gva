package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// CampingReservation 露营预约记录
type CampingReservation struct {
	global.GVA_MODEL
	SiteID      uint      `json:"siteId" form:"siteId" gorm:"column:site_id;comment:场地ID;"`
	ReserveDate time.Time `json:"reserveDate" form:"reserveDate" gorm:"column:reserve_date;comment:预约日期;type:date;"`
	TimeSlotID  uint      `json:"timeSlotId" form:"timeSlotId" gorm:"column:time_slot_id;comment:时段ID;"`
	BookerName  string    `json:"bookerName" form:"bookerName" gorm:"column:booker_name;comment:预订人;size:50;"`
	Phone       string    `json:"phone" form:"phone" gorm:"column:phone;comment:手机号;size:20;"`
	PeopleCount int       `json:"peopleCount" form:"peopleCount" gorm:"column:people_count;comment:人数;"`
	Remark      string    `json:"remark" form:"remark" gorm:"column:remark;comment:备注;type:text;"`
	VerifyCode  string    `json:"verifyCode" form:"verifyCode" gorm:"column:verify_code;comment:核销码;size:32;uniqueIndex;"`
	Status      int       `json:"status" form:"status" gorm:"column:status;comment:0待核销1已核销2已取消;default:0;"`
}

// TableName 表名
func (CampingReservation) TableName() string {
	return "gva_camping_reservations"
}
