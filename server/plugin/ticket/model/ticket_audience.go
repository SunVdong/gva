package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketAudience 门票适用人群
type TicketAudience struct {
	global.GVA_MODEL
	SkuID        uint   `json:"skuId" form:"skuId" gorm:"column:sku_id;comment:门票SKU ID;"`
	AudienceType string `json:"audienceType" form:"audienceType" gorm:"column:audience_type;comment:适用人群;size:50;"`
	Description  string `json:"description" form:"description" gorm:"column:description;comment:人群说明;size:255;"`
}

// TableName 表名
func (TicketAudience) TableName() string {
	return "ticket_audience"
}
