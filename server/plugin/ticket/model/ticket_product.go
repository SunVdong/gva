package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketProduct 门票商品表
type TicketProduct struct {
	global.GVA_MODEL
	ScenicID    uint   `json:"scenicId" form:"scenicId" gorm:"column:scenic_id;comment:景区ID;"`
	Name        string `json:"name" form:"name" gorm:"column:name;comment:门票商品名称;size:100;"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:门票说明;type:text;"`
	Status      int    `json:"status" form:"status" gorm:"column:status;comment:状态1启用0禁用;default:1;"`
	CreatedBy   int    `json:"createdBy" form:"createdBy" gorm:"column:created_by;default:0;"`
	UpdatedBy   int    `json:"updatedBy" form:"updatedBy" gorm:"column:updated_by;default:0;"`
}

// TableName 表名
func (TicketProduct) TableName() string {
	return "ticket_products"
}
