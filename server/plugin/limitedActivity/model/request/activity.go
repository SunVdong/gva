package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ActivitySearch 活动搜索
type ActivitySearch struct {
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
	request.PageInfo
}
