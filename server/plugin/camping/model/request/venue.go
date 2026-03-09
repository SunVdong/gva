package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// VenueSearch 场地搜索
type VenueSearch struct {
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
	request.PageInfo
}
