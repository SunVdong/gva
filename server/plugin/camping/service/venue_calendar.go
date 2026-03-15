package service

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
)

type venueCalendar struct{}

// dateOnly 截断为当日 0 点，避免 DATE 列带时分秒导致查询不匹配
func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (s *venueCalendar) SetOrCreate(venueID uint, date time.Time, status int) error {
	d := dateOnly(date)
	var cal model.VenueCalendar
	err := global.GVA_DB.Where("venue_id = ? AND date = ?", venueID, d).First(&cal).Error
	if err != nil {
		cal = model.VenueCalendar{VenueID: venueID, Date: d, Status: status}
		return global.GVA_DB.Create(&cal).Error
	}
	return global.GVA_DB.Model(&cal).Update("status", status).Error
}

func (s *venueCalendar) GetByVenueAndDateRange(venueID uint, start, end time.Time) (list []model.VenueCalendar, err error) {
	startD := dateOnly(start)
	endD := dateOnly(end)
	err = global.GVA_DB.Where("venue_id = ? AND date >= ? AND date <= ?", venueID, startD, endD).Find(&list).Error
	return
}

// IsDateOpen 某日是否可预约（无记录默认可约）
func (s *venueCalendar) IsDateOpen(venueID uint, date time.Time) (bool, error) {
	d := dateOnly(date)
	var cal model.VenueCalendar
	err := global.GVA_DB.Where("venue_id = ? AND date = ?", venueID, d).First(&cal).Error
	if err != nil {
		return true, nil
	}
	return cal.Status == 1, nil
}
