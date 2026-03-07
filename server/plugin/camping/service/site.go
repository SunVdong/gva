package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
)

type site struct{}

func (s *site) CreateSite(m *model.CampingSite) error {
	return global.GVA_DB.Create(m).Error
}

func (s *site) DeleteSite(id uint) error {
	return global.GVA_DB.Delete(&model.CampingSite{}, "id = ?", id).Error
}

func (s *site) DeleteSiteByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.CampingSite{}, "id in ?", ids).Error
}

func (s *site) UpdateSite(m model.CampingSite) error {
	return global.GVA_DB.Model(&model.CampingSite{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *site) GetSite(id uint) (model.CampingSite, error) {
	var res model.CampingSite
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *site) GetSiteList(req request.CampingSiteSearch) (list []model.CampingSite, total int64, err error) {
	db := global.GVA_DB.Model(&model.CampingSite{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
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
	err = db.Order("id DESC").Find(&list).Error
	return
}
