package model

import (
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// ActivityGuide 活动指南 结构体
type ActivityGuide struct {
	global.GVA_MODEL
	Name       string         `json:"name" form:"name" gorm:"column:name;comment:活动名称;size:128;"`                             // 活动名称
	Summary    string         `json:"summary" form:"summary" gorm:"column:summary;comment:简介;type:text;"`                     // 简介
	CoverImage string         `json:"coverImage" form:"coverImage" gorm:"column:cover_image;comment:封面图;size:512;"`         // 封面图
	Media      datatypes.JSON `json:"media" form:"media" gorm:"column:media;comment:介绍视频或图片;" swaggertype:"array,object"` // 介绍媒体 [{url,name}]，图片或仅 MP4 视频
	IsPreview  *bool          `json:"isPreview" form:"isPreview" gorm:"column:is_preview;comment:是否活动预告;default:false;"`   // 是否活动预告
	ShowStatus *bool          `json:"showStatus" form:"showStatus" gorm:"column:show_status;comment:显示状态;default:true;"`    // 显示状态
}

// TableName 活动指南自定义表名
func (ActivityGuide) TableName() string {
	return "gva_activity_guide"
}

type guideMediaItem struct {
	URL string `json:"url"`
}

// ValidateActivityGuideMedia 校验介绍媒体：支持常见图片后缀；视频仅允许 .mp4（忽略查询串）
func ValidateActivityGuideMedia(m datatypes.JSON) error {
	if len(m) == 0 {
		return nil
	}
	s := strings.TrimSpace(string(m))
	if s == "" || s == "null" {
		return nil
	}
	var items []guideMediaItem
	if err := json.Unmarshal(m, &items); err != nil {
		return errors.New("媒体数据格式错误")
	}
	for _, it := range items {
		u := strings.TrimSpace(it.URL)
		if u == "" {
			return errors.New("媒体项缺少地址")
		}
		if idx := strings.Index(u, "?"); idx >= 0 {
			u = u[:idx]
		}
		ext := strings.ToLower(path.Ext(u))
		if !isAllowedActivityGuideMediaExt(ext) {
			return errors.New("仅支持常见图片格式或 MP4 视频")
		}
	}
	return nil
}

func ValidateActivityGuideCoverImage(coverImage string) error {
	if strings.TrimSpace(coverImage) == "" {
		return errors.New("封面图必填")
	}
	return nil
}

func isAllowedActivityGuideMediaExt(ext string) bool {
	switch ext {
	case ".mp4", ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp":
		return true
	default:
		return false
	}
}
