package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ActivityOrderSearch 订单搜索
// OrderType: pending_payment|pending_verify|completed；不传返回全部
type ActivityOrderSearch struct {
	OrderNo      string  `json:"orderNo" form:"orderNo"`
	UserID       uint    `json:"userId" form:"userId"`
	ContactPhone string  `json:"contactPhone" form:"contactPhone"`
	ActivityID   uint    `json:"activityId" form:"activityId"`
	Status       *int    `json:"status" form:"status"`
	OrderType    *string `json:"orderType" form:"orderType"`
	request.PageInfo
}
