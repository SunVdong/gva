package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VenueReservation 场地预约订单
// 状态 0待核销 1已核销 2已取消 3已过期（核销后为已核销）
type VenueReservation struct {
	global.GVA_MODEL
	ReservationNo string     `json:"reservationNo" form:"reservationNo" gorm:"column:reservation_no;comment:预约单号;size:64;uniqueIndex;"`
	UserID        uint       `json:"userId" form:"userId" gorm:"column:user_id;comment:用户ID;"`
	VenueID       uint       `json:"venueId" form:"venueId" gorm:"column:venue_id;comment:场地ID;"`
	TimeslotID    uint       `json:"timeslotId" form:"timeslotId" gorm:"column:timeslot_id;comment:时间段ID;"`
	ReserveDate   time.Time  `json:"reserveDate" form:"reserveDate" gorm:"column:reserve_date;comment:预约日期;type:date;"`
	ContactName   string     `json:"contactName" form:"contactName" gorm:"column:contact_name;comment:联系人;size:50;"`
	ContactPhone  string     `json:"contactPhone" form:"contactPhone" gorm:"column:contact_phone;comment:联系电话;size:20;"`
	ContactCount  int        `json:"contactCount" form:"contactCount" gorm:"column:contact_count;comment:预约人数;"`
	Status        int        `json:"status" form:"status" gorm:"column:status;comment:0待核销1已核销2已取消3已过期;default:0;"`
	VerifyCode    string     `json:"verifyCode" form:"verifyCode" gorm:"column:verify_code;comment:核销码;size:32;index;"`
	VerifiedAt    *time.Time `json:"verifiedAt" form:"verifiedAt" gorm:"column:verified_at;comment:核销时间;"`
}

// TableName 表名
func (VenueReservation) TableName() string {
	return "venue_reservations"
}
