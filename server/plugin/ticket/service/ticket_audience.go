package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"gorm.io/gorm"
)

type ticketAudience struct{}

func (s *ticketAudience) GetBySku(skuId uint) ([]model.TicketAudience, error) {
	var list []model.TicketAudience
	err := global.GVA_DB.Where("sku_id = ?", skuId).Find(&list).Error
	return list, err
}

func (s *ticketAudience) SaveBySku(skuId uint, items []model.TicketAudience) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("sku_id = ?", skuId).Delete(&model.TicketAudience{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].SkuID = skuId
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
