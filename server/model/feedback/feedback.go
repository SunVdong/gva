package feedback

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
)

// Feedback 意见反馈（须登录，user_id 为 C 端 users 表主键）
type Feedback struct {
	global.GVA_MODEL
	UserID  uint      `json:"userId" gorm:"column:user_id;index;not null;comment:C端用户ID"`
	Content string    `json:"content" gorm:"column:content;type:text;comment:反馈内容"`
	User    user.User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}
