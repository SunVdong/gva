package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// TicketProductSearch 门票商品搜索
type TicketProductSearch struct {
	ScenicID uint   `json:"scenicId" form:"scenicId"`
	Name     string `json:"name" form:"name"`
	Status   *int   `json:"status" form:"status"`
	request.PageInfo
}
