package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// BannerSearch Banner 搜索
type BannerSearch struct {
	Title  string `json:"title" form:"title"`
	Status *int   `json:"status" form:"status"`
	request.PageInfo
}
