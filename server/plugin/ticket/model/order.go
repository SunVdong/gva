package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TicketOrder 订单表（一单一SKU，原 order_items 字段已合并）
// 状态：0待支付 1待核销 2已核销 3已取消 4已过期 5已关闭（核销后 VerifiedAt 不为空，仅已核销订单可评价）
type TicketOrder struct {
	global.GVA_MODEL
	OrderNo     string     `json:"orderNo" gorm:"column:order_no;comment:订单号;size:64;uniqueIndex;"`
	UserID      uint       `json:"userId" gorm:"column:user_id;comment:用户ID;"`
	BookerName  string     `json:"bookerName" gorm:"column:booker_name;comment:预定人姓名;size:50;"`
	BookerPhone string     `json:"bookerPhone" gorm:"column:booker_phone;comment:预定人手机号;size:20;"`
	SkuID       uint       `json:"skuId" gorm:"column:sku_id;comment:门票SKU ID;"`
	SkuName     string     `json:"skuName" gorm:"column:sku_name;comment:SKU名称;size:50;"`
	Price       float64    `json:"price" gorm:"column:price;type:decimal(10,2);comment:购买单价;"`
	Quantity    int        `json:"quantity" gorm:"column:quantity;comment:购买数量;"`
	VisitDate   time.Time  `json:"visitDate" gorm:"column:visit_date;type:date;comment:游玩日期;"`
	TotalAmount float64    `json:"totalAmount" gorm:"column:total_amount;type:decimal(10,2);comment:订单总金额;"`
	PayAmount   float64    `json:"payAmount" gorm:"column:pay_amount;type:decimal(10,2);comment:支付金额;"`
	Status        int        `json:"status" gorm:"column:status;comment:订单状态0待支付1待核销2已核销3已取消4已过期5已关闭;default:0;"`
	TotalUseTimes int        `json:"totalUseTimes" gorm:"column:total_use_times;default:0;comment:总可核销次数;"`
	VerifiedTimes int        `json:"verifiedTimes" gorm:"column:verified_times;default:0;comment:已核销次数;"`
	PayTime       *time.Time `json:"payTime" gorm:"column:pay_time;comment:支付时间;"`
	VerifiedAt    *time.Time `json:"verifiedAt" gorm:"column:verified_at;comment:核销时间;"`
	ProductName string     `json:"productName" gorm:"-"`
}

// TableName 表名
func (TicketOrder) TableName() string {
	return "orders"
}
