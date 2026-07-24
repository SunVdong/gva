package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Banner 活动列表轮播
// 显示状态：1 显示 / 0 隐藏
type Banner struct {
	global.GVA_MODEL
	Title       string `json:"title" form:"title" gorm:"column:title;comment:后台识别标题;size:128;"`
	Image       string `json:"image" form:"image" gorm:"column:image;comment:轮播图;size:512;"`
	DetailImage string `json:"detailImage" form:"detailImage" gorm:"column:detail_image;comment:详情长图;size:512;"`
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序越小越靠前;default:0;"`
	Status      int    `json:"status" form:"status" gorm:"column:status;comment:显示状态1显示0隐藏;default:1;"`
	CreatedBy   int    `json:"createdBy" form:"createdBy" gorm:"column:created_by;default:0;"`
	UpdatedBy   int    `json:"updatedBy" form:"updatedBy" gorm:"column:updated_by;default:0;"`
}

// TableName 表名
func (Banner) TableName() string {
	return "limited_activity_banners"
}
