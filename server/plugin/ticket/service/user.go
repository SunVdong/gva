package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
)

type ticketUser struct{}

func (s *ticketUser) GetList(req request.TicketUserSearch) (list []model.TicketUser, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketUser{})
	if req.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Openid != "" {
		db = db.Where("openid LIKE ?", "%"+req.Openid+"%")
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

func (s *ticketUser) Get(id uint) (model.TicketUser, error) {
	var res model.TicketUser
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *ticketUser) Update(m model.TicketUser) error {
	return global.GVA_DB.Model(&model.TicketUser{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"nickname": m.Nickname,
		"avatar":   m.Avatar,
		"phone":    m.Phone,
	}).Error
}
