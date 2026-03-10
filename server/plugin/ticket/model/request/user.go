package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// TicketUserSearch 用户搜索
type TicketUserSearch struct {
	Nickname string `json:"nickname" form:"nickname"`
	Phone    string `json:"phone" form:"phone"`
	Openid   string `json:"openid" form:"openid"`
	request.PageInfo
}
