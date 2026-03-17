package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// TicketOrderSearch 订单搜索
// OrderType: 待支付=pending_payment, 待核销=pending_verify, 已完成=completed；不传返回全部
type TicketOrderSearch struct {
	OrderNo   string  `json:"orderNo" form:"orderNo"`
	UserID    uint    `json:"userId" form:"userId"`
	Status    *int    `json:"status" form:"status"`       // 0待支付 1已支付 2已退款（后台用）
	OrderType *string `json:"orderType" form:"orderType"` // 待支付/待核销/已完成，不传默认全部
	request.PageInfo
}
