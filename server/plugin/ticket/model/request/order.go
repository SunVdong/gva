package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// TicketOrderSearch 订单搜索
type TicketOrderSearch struct {
	OrderNo string `json:"orderNo" form:"orderNo"`
	UserID  uint   `json:"userId" form:"userId"`
	Status  *int   `json:"status" form:"status"` // 0待支付 1已支付 2已退款
	request.PageInfo
}
