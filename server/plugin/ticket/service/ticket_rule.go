package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"gorm.io/gorm"
)

type ticketRule struct{}

func (s *ticketRule) GetByProduct(productId uint) ([]model.TicketRule, error) {
	var list []model.TicketRule
	err := global.GVA_DB.Where("product_id = ?", productId).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (s *ticketRule) SaveByProduct(productId uint, list []model.TicketRule) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productId).Delete(&model.TicketRule{}).Error; err != nil {
			return err
		}
		for i := range list {
			list[i].ProductID = productId
			if err := tx.Create(&list[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
