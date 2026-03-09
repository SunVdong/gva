package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
)

type venueOpenTime struct{}

func (s *venueOpenTime) CreateVenueOpenTime(m *model.VenueOpenTime) error {
	return global.GVA_DB.Create(m).Error
}

func (s *venueOpenTime) DeleteVenueOpenTime(id uint) error {
	return global.GVA_DB.Delete(&model.VenueOpenTime{}, "id = ?", id).Error
}

func (s *venueOpenTime) UpdateVenueOpenTime(m model.VenueOpenTime) error {
	return global.GVA_DB.Model(&model.VenueOpenTime{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *venueOpenTime) GetVenueOpenTime(id uint) (model.VenueOpenTime, error) {
	var res model.VenueOpenTime
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *venueOpenTime) GetVenueOpenTimeListByVenue(venueID uint) (list []model.VenueOpenTime, err error) {
	err = global.GVA_DB.Where("venue_id = ?", venueID).Order("week_day ASC").Find(&list).Error
	return
}

func (s *venueOpenTime) SaveVenueOpenTimes(venueID uint, list []request.VenueOpenTimeBody) error {
	if err := global.GVA_DB.Where("venue_id = ?", venueID).Delete(&model.VenueOpenTime{}).Error; err != nil {
		return err
	}
	for _, item := range list {
		m := model.VenueOpenTime{
			VenueID:   venueID,
			WeekDay:   item.WeekDay,
			OpenTime:  model.TimeOnly(item.OpenTime),
			CloseTime: model.TimeOnly(item.CloseTime),
		}
		if err := global.GVA_DB.Create(&m).Error; err != nil {
			return err
		}
	}
	return nil
}
