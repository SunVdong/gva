package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketOrder 订单表
type TicketOrder struct {
	global.GVA_MODEL
	OrderNo     string     `json:"orderNo" gorm:"column:order_no;comment:订单号;size:64;uniqueIndex;"`
	UserID      uint       `json:"userId" gorm:"column:user_id;comment:用户ID;"`
	TotalAmount float64    `json:"totalAmount" gorm:"column:total_amount;type:decimal(10,2);comment:订单总金额;"`
	PayAmount   float64    `json:"payAmount" gorm:"column:pay_amount;type:decimal(10,2);comment:支付金额;"`
	Status      int        `json:"status" gorm:"column:status;comment:订单状态0待支付1已支付2已退款;default:0;"`
	PayTime     *time.Time `json:"payTime" gorm:"column:pay_time;comment:支付时间;"`
}

// TableName 表名
func (TicketOrder) TableName() string {
	return "orders"
}
