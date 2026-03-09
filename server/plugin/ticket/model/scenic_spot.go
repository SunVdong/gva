package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// ScenicSpot 景区表
type ScenicSpot struct {
	global.GVA_MODEL
	Name              string         `json:"name" form:"name" gorm:"column:name;comment:景区名称;size:100;"`
	CarouselImages    datatypes.JSON `json:"carouselImages" form:"carouselImages" gorm:"column:carousel_images;comment:轮播图URL列表;type:json;"`
	Description       string         `json:"description" form:"description" gorm:"column:description;comment:景区介绍;type:text;"`
	RefundChangeHours int            `json:"refundChangeHours" form:"refundChangeHours" gorm:"column:refund_change_hours;comment:可退订时间(小时)0不可退订;default:0;"`
	Status            int            `json:"status" form:"status" gorm:"column:status;comment:状态1启用0禁用;default:1;"`
	CreatedBy         int            `json:"createdBy" form:"createdBy" gorm:"column:created_by;default:0;"`
	UpdatedBy         int            `json:"updatedBy" form:"updatedBy" gorm:"column:updated_by;default:0;"`
}

// TableName 表名
func (ScenicSpot) TableName() string {
	return "scenic_spots"
}
