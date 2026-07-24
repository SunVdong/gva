package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ActivityOrderVerifyRecord 活动订单核销记录
type ActivityOrderVerifyRecord struct {
	global.GVA_MODEL
	OrderID    uint      `json:"orderId" gorm:"column:order_id;comment:订单ID;index;"`
	VerifyNo   int       `json:"verifyNo" gorm:"column:verify_no;comment:第几次核销;"`
	VerifiedAt time.Time `json:"verifiedAt" gorm:"column:verified_at;comment:核销时间;"`
	Remark     string    `json:"remark" gorm:"column:remark;comment:备注;size:200;"`
}

// TableName 表名
func (ActivityOrderVerifyRecord) TableName() string {
	return "limited_activity_order_verify_records"
}
