package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"gorm.io/gorm"
)

type scenicOpenTime struct{}

func (s *scenicOpenTime) GetByScenic(scenicId uint) ([]model.ScenicOpenTime, error) {
	var list []model.ScenicOpenTime
	err := global.GVA_DB.Where("scenic_id = ?", scenicId).Order("week_day").Find(&list).Error
	return list, err
}

func (s *scenicOpenTime) Save(req request.ScenicOpenTimeSave) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("scenic_id = ?", req.ScenicID).Delete(&model.ScenicOpenTime{}).Error; err != nil {
			return err
		}
		for _, item := range req.List {
			m := model.ScenicOpenTime{
				ScenicID:  req.ScenicID,
				WeekDay:   item.WeekDay,
				OpenTime:  model.TimeOnly(item.OpenTime),
				CloseTime: model.TimeOnly(item.CloseTime),
			}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
