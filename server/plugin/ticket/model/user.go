package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketUser C 端用户表（如小程序用户）
type TicketUser struct {
	global.GVA_MODEL
	Openid   string `json:"openid" gorm:"column:openid;comment:微信openid;size:64;uniqueIndex;"`
	Nickname string `json:"nickname" gorm:"column:nickname;comment:昵称;size:100;"`
	Avatar   string `json:"avatar" gorm:"column:avatar;comment:头像;size:255;"`
	Phone    string `json:"phone" gorm:"column:phone;comment:手机号;size:20;"`
}

// TableName 表名
func (TicketUser) TableName() string {
	return "users"
}
