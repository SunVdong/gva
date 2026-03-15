package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	userModel "github.com/flipped-aurora/gin-vue-admin/server/model/user"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
)

type ticketUser struct{}

func (s *ticketUser) GetList(req request.TicketUserSearch) (list []userModel.User, total int64, err error) {
	db := global.GVA_DB.Model(&userModel.User{})
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

func (s *ticketUser) Get(id uint) (userModel.User, error) {
	var res userModel.User
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *ticketUser) Update(req request.TicketUserUpdate) error {
	updates := map[string]interface{}{
		"nickname":   req.Nickname,
		"avatar_url": req.Avatar,
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	} else {
		updates["phone"] = nil
	}
	return global.GVA_DB.Model(&userModel.User{}).Where("id = ?", req.ID).Updates(updates).Error
}
