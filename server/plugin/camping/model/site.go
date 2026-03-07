package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// CampingSite 露营场地
type CampingSite struct {
	global.GVA_MODEL
	Name           string         `json:"name" form:"name" gorm:"column:name;comment:场地名称;size:100;"`
	CarouselImages datatypes.JSON `json:"carouselImages" form:"carouselImages" gorm:"column:carousel_images;comment:轮播图URL列表;type:json;"`
	Introduction   string         `json:"introduction" form:"introduction" gorm:"column:introduction;comment:场地介绍富文本;type:longtext;"`
	ReserveRules   string         `json:"reserveRules" form:"reserveRules" gorm:"column:reserve_rules;comment:预约规则介绍富文本;type:longtext;"`
	OpenTimeDesc   string         `json:"openTimeDesc" form:"openTimeDesc" gorm:"column:open_time_desc;comment:开放时间说明;type:text;"`
	Status         int            `json:"status" form:"status" gorm:"column:status;comment:状态0禁用1启用;default:1;"`
}

// TableName 表名
func (CampingSite) TableName() string {
	return "gva_camping_sites"
}
