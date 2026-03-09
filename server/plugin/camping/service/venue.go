package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
)

type venue struct{}

func (s *venue) CreateVenue(m *model.Venue) error {
	return global.GVA_DB.Create(m).Error
}

func (s *venue) DeleteVenue(id uint) error {
	return global.GVA_DB.Delete(&model.Venue{}, "id = ?", id).Error
}

func (s *venue) DeleteVenueByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.Venue{}, "id in ?", ids).Error
}

func (s *venue) UpdateVenue(m model.Venue) error {
	return global.GVA_DB.Model(&model.Venue{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *venue) GetVenue(id uint) (model.Venue, error) {
	var res model.Venue
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *venue) GetVenueList(req request.VenueSearch) (list []model.Venue, total int64, err error) {
	db := global.GVA_DB.Model(&model.Venue{})
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
