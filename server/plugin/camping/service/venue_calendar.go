package service

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
)

type venueCalendar struct{}

func (s *venueCalendar) SetOrCreate(venueID uint, date time.Time, status int) error {
	var cal model.VenueCalendar
	err := global.GVA_DB.Where("venue_id = ? AND date = ?", venueID, date).First(&cal).Error
	if err != nil {
		cal = model.VenueCalendar{VenueID: venueID, Date: date, Status: status}
		return global.GVA_DB.Create(&cal).Error
	}
	return global.GVA_DB.Model(&cal).Update("status", status).Error
}

func (s *venueCalendar) GetByVenueAndDateRange(venueID uint, start, end time.Time) (list []model.VenueCalendar, err error) {
	err = global.GVA_DB.Where("venue_id = ? AND date >= ? AND date <= ?", venueID, start, end).Find(&list).Error
	return
}

// IsDateOpen 某日是否可预约（无记录默认可约）
func (s *venueCalendar) IsDateOpen(venueID uint, date time.Time) (bool, error) {
	var cal model.VenueCalendar
	err := global.GVA_DB.Where("venue_id = ? AND date = ?", venueID, date).First(&cal).Error
	if err != nil {
		return true, nil
	}
	return cal.Status == 1, nil
}
