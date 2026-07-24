package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// OrderVerifyRecord 订单核销记录表
type OrderVerifyRecord struct {
	global.GVA_MODEL
	OrderID    uint      `json:"orderId" gorm:"column:order_id;comment:订单ID;index;"`
	VerifyNo   int       `json:"verifyNo" gorm:"column:verify_no;comment:第几次核销;"`
	VerifiedAt time.Time `json:"verifiedAt" gorm:"column:verified_at;comment:核销时间;"`
	Venue      string    `json:"venue" gorm:"column:venue;size:32;comment:核销场合code不支持多场合时为空;"`
	Remark     string    `json:"remark" gorm:"column:remark;comment:备注;size:200;"`
}

// TableName 表名
func (OrderVerifyRecord) TableName() string {
	return "order_verify_records"
}
