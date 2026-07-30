package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// TicketOrderSearch 订单搜索
// OrderType: 待支付=pending_payment, 待核销=pending_verify, 已完成=completed；不传返回全部
type TicketOrderSearch struct {
	OrderNo     string  `json:"orderNo" form:"orderNo"`
	UserID      uint    `json:"userId" form:"userId"`
	BookerPhone string  `json:"bookerPhone" form:"bookerPhone"`
	Status      *int    `json:"status" form:"status"`       // 0待支付 1待核销 2已核销 3已取消 4已过期 5已关闭 6已退款 7退款中（后台用）
	OrderType   *string `json:"orderType" form:"orderType"` // 待支付/待核销/已完成，不传默认全部
	TicketType  *int    `json:"ticketType" form:"ticketType"` // 票种：1单次 2多次（join 当前 SKU）
	Venue       string  `json:"venue" form:"venue"`           // 核销场合 code：至少一次匹配
	request.PageInfo
}

// TicketVenueVerifyStatsReq 核销场合按月统计
type TicketVenueVerifyStatsReq struct {
	Month string `json:"month" form:"month"` // YYYY-MM，空则当前月
}

// TicketVenueVerifyStatsItem 单场合核销次数
type TicketVenueVerifyStatsItem struct {
	Venue string `json:"venue"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}
