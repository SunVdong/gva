package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// OrderItem 订单项表
type OrderItem struct {
	global.GVA_MODEL
	OrderID   uint      `json:"orderId" gorm:"column:order_id;comment:订单ID;"`
	SkuID     uint      `json:"skuId" gorm:"column:sku_id;comment:门票SKU ID;"`
	SkuName   string    `json:"skuName" gorm:"column:sku_name;comment:SKU名称;size:50;"`
	Price     float64   `json:"price" gorm:"column:price;type:decimal(10,2);comment:购买价格;"`
	Quantity  int       `json:"quantity" gorm:"column:quantity;comment:购买数量;"`
	VisitDate time.Time  `json:"visitDate" gorm:"column:visit_date;type:date;comment:游玩日期;"`
}

// TableName 表名
func (OrderItem) TableName() string {
	return "order_items"
}
