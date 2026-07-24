package service

import (
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
)

type banner struct{}

func (s *banner) validate(m *model.Banner) error {
	m.Title = strings.TrimSpace(m.Title)
	m.Image = strings.TrimSpace(m.Image)
	m.DetailImage = strings.TrimSpace(m.DetailImage)
	if m.Image == "" {
		return fmt.Errorf("轮播图不能为空")
	}
	if m.DetailImage == "" {
		return fmt.Errorf("详情长图不能为空")
	}
	return nil
}

func (s *banner) Create(m *model.Banner) error {
	if err := s.validate(m); err != nil {
		return err
	}
	return global.GVA_DB.Create(m).Error
}

func (s *banner) Delete(id uint) error {
	return global.GVA_DB.Delete(&model.Banner{}, "id = ?", id).Error
}

func (s *banner) DeleteByIds(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.GVA_DB.Delete(&model.Banner{}, "id IN ?", ids).Error
}

func (s *banner) Update(m model.Banner) error {
	var old model.Banner
	if err := global.GVA_DB.Where("id = ?", m.ID).First(&old).Error; err != nil {
		return fmt.Errorf("Banner 不存在")
	}
	if err := s.validate(&m); err != nil {
		return err
	}
	updates := map[string]any{
		"title":        m.Title,
		"image":        m.Image,
		"detail_image": m.DetailImage,
		"sort":         m.Sort,
		"status":       m.Status,
		"updated_by":   m.UpdatedBy,
	}
	return global.GVA_DB.Model(&model.Banner{}).Where("id = ?", m.ID).Updates(updates).Error
}

func (s *banner) Get(id uint) (model.Banner, error) {
	var res model.Banner
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *banner) GetList(req request.BannerSearch) (list []model.Banner, total int64, err error) {
	db := global.GVA_DB.Model(&model.Banner{})
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("sort ASC, id ASC").Find(&list).Error
	return
}

// GetMiniList 小程序启用中的 Banner（status=1，按 sort ASC, id ASC）
func (s *banner) GetMiniList() (list []model.Banner, err error) {
	err = global.GVA_DB.Model(&model.Banner{}).
		Where("status = ?", 1).
		Order("sort ASC, id ASC").
		Find(&list).Error
	return
}
