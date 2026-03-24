package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// TicketSkuSearch 门票SKU搜索
type TicketSkuSearch struct {
	ProductID uint `json:"productId" form:"productId"`
	Status    *int `json:"status" form:"status"`
	request.PageInfo
}