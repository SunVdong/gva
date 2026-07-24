package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ActivityOrder 限时活动参与订单
// 状态：0待支付 1待核销 2已核销 5已关闭 6已退款 7退款中
type ActivityOrder struct {
	global.GVA_MODEL
	OrderNo         string     `json:"orderNo" gorm:"column:order_no;comment:订单号;size:64;uniqueIndex;"`
	UserID          uint       `json:"userId" gorm:"column:user_id;comment:用户ID;index;"`
	ActivityID      uint       `json:"activityId" gorm:"column:activity_id;comment:活动ID;index;"`
	ActivityName    string     `json:"activityName" gorm:"column:activity_name;comment:活动名称快照;size:128;"`
	ContactName     string     `json:"contactName" gorm:"column:contact_name;comment:联系人姓名;size:50;"`
	ContactPhone    string     `json:"contactPhone" gorm:"column:contact_phone;comment:联系人手机号;size:20;"`
	Quantity        int        `json:"quantity" gorm:"column:quantity;comment:人次;"`
	UnitPrice       float64    `json:"unitPrice" gorm:"column:unit_price;type:decimal(10,2);comment:实际单价;"`
	PayAmount       float64    `json:"payAmount" gorm:"column:pay_amount;type:decimal(10,2);comment:应付总额;"`
	Status          int        `json:"status" gorm:"column:status;comment:订单状态0待支付1待核销2已核销5已关闭6已退款7退款中;default:0;"`
	TotalUseTimes   int        `json:"totalUseTimes" gorm:"column:total_use_times;default:0;comment:总可核销次数;"`
	VerifiedTimes   int        `json:"verifiedTimes" gorm:"column:verified_times;default:0;comment:已核销次数;"`
	PayTime         *time.Time `json:"payTime" gorm:"column:pay_time;comment:支付时间;"`
	WxTransactionID string     `json:"wxTransactionId" gorm:"column:wx_transaction_id;size:64;index;comment:微信支付订单号;"`
	RefundNo        string     `json:"refundNo" gorm:"column:refund_no;size:64;index;comment:商户退款单号;"`
	WxRefundID      string     `json:"wxRefundId" gorm:"column:wx_refund_id;size:64;comment:微信退款单号;"`
	RefundTime      *time.Time `json:"refundTime" gorm:"column:refund_time;comment:退款时间;"`
	RefundAmount    float64    `json:"refundAmount" gorm:"column:refund_amount;type:decimal(10,2);comment:实退金额(元);default:0;"`
	VerifiedAt      *time.Time `json:"verifiedAt" gorm:"column:verified_at;comment:全部核销完成时间;"`
	// 虚拟：活动群二维码 / 客服二维码（详情填充）
	GroupQr   string `json:"groupQr" gorm:"-"`
	ServiceQr string `json:"serviceQr" gorm:"-"`
	// 虚拟：活动封面
	CoverImage string `json:"coverImage" gorm:"-"`
}

// TableName 表名
func (ActivityOrder) TableName() string {
	return "limited_activity_orders"
}
