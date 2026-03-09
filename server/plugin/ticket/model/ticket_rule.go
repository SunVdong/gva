package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketRule 门票规则表
type TicketRule struct {
	global.GVA_MODEL
	ProductID uint   `json:"productId" form:"productId" gorm:"column:product_id;comment:门票商品ID;"`
	Title     string `json:"title" form:"title" gorm:"column:title;comment:规则标题;size:100;"`
	Content   string `json:"content" form:"content" gorm:"column:content;comment:规则内容;type:text;"`
	Sort      int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`
}

// TableName 表名
func (TicketRule) TableName() string {
	return "ticket_rules"
}
