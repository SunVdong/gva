package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketSku 门票SKU表
type TicketSku struct {
	global.GVA_MODEL
	ProductID   uint    `json:"productId" form:"productId" gorm:"column:product_id;comment:门票商品ID;"`
	Name        string  `json:"name" form:"name" gorm:"column:name;comment:SKU名称如成人票;size:50;"`
	Price       float64 `json:"price" form:"price" gorm:"column:price;type:decimal(10,2);comment:销售价格;"`
	MarketPrice *float64 `json:"marketPrice" form:"marketPrice" gorm:"column:market_price;type:decimal(10,2);comment:市场价格;"`
	Stock       int     `json:"stock" form:"stock" gorm:"column:stock;default:0;comment:库存;"`
	LimitBuy    int     `json:"limitBuy" form:"limitBuy" gorm:"column:limit_buy;default:0;comment:每人限购数量;"`
	Sort          int     `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序值越小越靠前;"`
	Status        int     `json:"status" form:"status" gorm:"column:status;comment:状态1启用0禁用;default:1;"`
	BookingNotice string  `json:"bookingNotice" form:"bookingNotice" gorm:"column:booking_notice;type:text;comment:预定须知;"`
	CreatedBy   int     `json:"createdBy" form:"createdBy" gorm:"column:created_by;default:0;"`
	UpdatedBy   int     `json:"updatedBy" form:"updatedBy" gorm:"column:updated_by;default:0;"`
}

// TableName 表名
func (TicketSku) TableName() string {
	return "ticket_sku"
}
