package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketCalendar 门票日历库存
type TicketCalendar struct {
	global.GVA_MODEL
	SkuID     uint      `json:"skuId" form:"skuId" gorm:"column:sku_id;comment:门票SKU ID;uniqueIndex:idx_ticket_calendar_sku_date;"`
	VisitDate time.Time `json:"visitDate" form:"visitDate" gorm:"column:visit_date;type:date;comment:游玩日期;uniqueIndex:idx_ticket_calendar_sku_date;"`
	Stock     int       `json:"stock" form:"stock" gorm:"column:stock;default:0;comment:库存;"`
	Sold      int       `json:"sold" form:"sold" gorm:"column:sold;default:0;comment:已售数量;"`
	Status    int       `json:"status" form:"status" gorm:"column:status;comment:状态1可售0关闭;default:1;"`
}

// TableName 表名
func (TicketCalendar) TableName() string {
	return "ticket_calendar"
}
