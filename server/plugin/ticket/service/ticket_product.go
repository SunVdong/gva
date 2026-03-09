package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
)

type ticketProduct struct{}

func (s *ticketProduct) Create(m *model.TicketProduct) error {
	return global.GVA_DB.Create(m).Error
}

func (s *ticketProduct) Delete(id uint) error {
	return global.GVA_DB.Delete(&model.TicketProduct{}, "id = ?", id).Error
}

func (s *ticketProduct) DeleteByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.TicketProduct{}, "id in ?", ids).Error
}

func (s *ticketProduct) Update(m model.TicketProduct) error {
	return global.GVA_DB.Model(&model.TicketProduct{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *ticketProduct) Get(id uint) (model.TicketProduct, error) {
	var res model.TicketProduct
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *ticketProduct) GetList(req request.TicketProductSearch) (list []model.TicketProduct, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketProduct{})
	if req.ScenicID > 0 {
		db = db.Where("scenic_id = ?", req.ScenicID)
	}
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
