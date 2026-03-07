package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CampingSiteSearch 场地搜索
type CampingSiteSearch struct {
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
	request.PageInfo
}
