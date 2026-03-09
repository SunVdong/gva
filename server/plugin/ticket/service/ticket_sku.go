package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
)

type ticketSku struct{}

func (s *ticketSku) Create(m *model.TicketSku) error {
	return global.GVA_DB.Create(m).Error
}

func (s *ticketSku) Delete(id uint) error {
	return global.GVA_DB.Delete(&model.TicketSku{}, "id = ?", id).Error
}

func (s *ticketSku) DeleteByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.TicketSku{}, "id in ?", ids).Error
}

func (s *ticketSku) Update(m model.TicketSku) error {
	return global.GVA_DB.Model(&model.TicketSku{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *ticketSku) Get(id uint) (model.TicketSku, error) {
	var res model.TicketSku
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *ticketSku) GetList(req request.TicketSkuSearch) (list []model.TicketSku, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketSku{})
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
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
