package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Activity 限时活动
// 显示状态：1 显示 / 0 隐藏
type Activity struct {
	global.GVA_MODEL
	Name        string    `json:"name" form:"name" gorm:"column:name;comment:活动名称;size:128;"`
	Detail      string    `json:"detail" form:"detail" gorm:"column:detail;comment:活动详情;type:text;"`
	StartTime   time.Time `json:"startTime" form:"startTime" gorm:"column:start_time;comment:活动开始时间;index;"`
	EndTime     time.Time `json:"endTime" form:"endTime" gorm:"column:end_time;comment:活动结束时间;index;"`
	MarketPrice float64   `json:"marketPrice" form:"marketPrice" gorm:"column:market_price;type:decimal(10,2);comment:市场价;"`
	Price       float64   `json:"price" form:"price" gorm:"column:price;type:decimal(10,2);comment:实际价格;"`
	Quota       int       `json:"quota" form:"quota" gorm:"column:quota;comment:总名额(人次);default:0;"`
	Sold        int       `json:"sold" form:"sold" gorm:"column:sold;comment:已占用名额(含待支付);default:0;"`
	CoverImage  string    `json:"coverImage" form:"coverImage" gorm:"column:cover_image;comment:封面图;size:512;"`
	LongImage   string    `json:"longImage" form:"longImage" gorm:"column:long_image;comment:长图(点击封面跳转);size:512;"`
	GroupQr     string    `json:"groupQr" form:"groupQr" gorm:"column:group_qr;comment:群二维码;size:512;"`
	ServiceQr   string    `json:"serviceQr" form:"serviceQr" gorm:"column:service_qr;comment:客服二维码;size:512;"`
	Status      int       `json:"status" form:"status" gorm:"column:status;comment:显示状态1显示0隐藏;default:1;"`
	CreatedBy   int       `json:"createdBy" form:"createdBy" gorm:"column:created_by;default:0;"`
	UpdatedBy   int       `json:"updatedBy" form:"updatedBy" gorm:"column:updated_by;default:0;"`
	// 虚拟字段：剩余可报名名额
	Remaining int `json:"remaining" gorm:"-"`
	// 虚拟字段：当前是否可报名（显示中且在时间窗内且有余量）
	CanSignup bool `json:"canSignup" gorm:"-"`
}

// TableName 表名
func (Activity) TableName() string {
	return "limited_activities"
}
