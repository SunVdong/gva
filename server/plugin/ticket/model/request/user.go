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

// TicketUserUpdate 后台更新 C 端用户（与 server/model/user.User 对应表 users）
type TicketUserUpdate struct {
	ID       uint   `json:"id" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`   // 对应 User.AvatarURL 列 avatar_url
	Phone    string `json:"phone"`
}
