package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// ActivityGuide 活动指南 结构体
type ActivityGuide struct {
	global.GVA_MODEL
	Name       string         `json:"name" form:"name" gorm:"column:name;comment:活动名称;size:128;"`                    // 活动名称
	Summary    string         `json:"summary" form:"summary" gorm:"column:summary;comment:简介;type:text;"`            // 简介
	CoverImage string         `json:"coverImage" form:"coverImage" gorm:"column:cover_image;comment:封面图;size:512;"`    // 封面图
	Media      datatypes.JSON `json:"media" form:"media" gorm:"column:media;comment:介绍视频或图片;" swaggertype:"array,object"` // 介绍视频或图片 [{type,url,name}]
	ShowStatus *bool          `json:"showStatus" form:"showStatus" gorm:"column:show_status;comment:显示状态;default:true;"`   // 显示状态
}

// TableName 活动指南自定义表名
func (ActivityGuide) TableName() string {
	return "gva_activity_guide"
}
