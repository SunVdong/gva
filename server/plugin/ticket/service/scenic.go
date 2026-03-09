package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
)

type scenic struct{}

func (s *scenic) Create(m *model.ScenicSpot) error {
	return global.GVA_DB.Create(m).Error
}

func (s *scenic) Delete(id uint) error {
	return global.GVA_DB.Delete(&model.ScenicSpot{}, "id = ?", id).Error
}

func (s *scenic) DeleteByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.ScenicSpot{}, "id in ?", ids).Error
}

func (s *scenic) Update(m model.ScenicSpot) error {
	return global.GVA_DB.Model(&model.ScenicSpot{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *scenic) Get(id uint) (model.ScenicSpot, error) {
	var res model.ScenicSpot
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *scenic) GetList(req request.ScenicSearch) (list []model.ScenicSpot, total int64, err error) {
	db := global.GVA_DB.Model(&model.ScenicSpot{})
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
