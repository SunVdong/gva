package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ScenicSearch 景区搜索
type ScenicSearch struct {
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
	request.PageInfo
}
