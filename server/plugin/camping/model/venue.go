package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// Venue 场地表（露营区、VIP场地、烧烤区等）
type Venue struct {
	global.GVA_MODEL
	Name              string         `json:"name" form:"name" gorm:"column:name;comment:场地名称;size:100;"`
	Description       string         `json:"description" form:"description" gorm:"column:description;comment:场地介绍富文本;type:longtext;"`
	CarouselImages    datatypes.JSON `json:"carouselImages" form:"carouselImages" gorm:"column:carousel_images;comment:轮播图URL列表;type:json;"`
	ReserveRules      string         `json:"reserveRules" form:"reserveRules" gorm:"column:reserve_rules;comment:预约规则富文本;type:longtext;"`
	RefundChangeHours int            `json:"refundChangeHours" form:"refundChangeHours" gorm:"column:refund_change_hours;comment:可退改时间(小时)0不可退改;default:0;"`
	Status            int            `json:"status" form:"status" gorm:"column:status;comment:状态1可用0关闭;default:1;"`
}

// TableName 表名
func (Venue) TableName() string {
	return "venues"
}
