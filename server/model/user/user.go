package user

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// User C 端/小程序用户（与后台 sys_users 分离），表名 users
type User struct {
	global.GVA_MODEL
	OpenID     string `json:"openid" gorm:"column:openid;uniqueIndex;comment:微信openid;size:64"`
	UnionID    string `json:"unionid" gorm:"column:unionid;index;comment:微信unionid;size:64"`
	SessionKey string `json:"-" gorm:"column:session_key;comment:会话密钥;size:64"`
	Nickname   string `json:"nickname" gorm:"column:nickname;comment:昵称;size:64"`
	AvatarURL  string `json:"avatarUrl" gorm:"column:avatar_url;comment:头像;size:255"`
}

func (User) TableName() string {
	return "users"
}
