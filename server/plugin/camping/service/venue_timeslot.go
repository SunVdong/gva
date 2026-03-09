package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
)

type venueTimeslot struct{}

func (s *venueTimeslot) CreateVenueTimeslot(m *model.VenueTimeslot) error {
	return global.GVA_DB.Create(m).Error
}

func (s *venueTimeslot) DeleteVenueTimeslot(id uint) error {
	return global.GVA_DB.Delete(&model.VenueTimeslot{}, "id = ?", id).Error
}

func (s *venueTimeslot) DeleteVenueTimeslotByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.VenueTimeslot{}, "id in ?", ids).Error
}

func (s *venueTimeslot) UpdateVenueTimeslot(m model.VenueTimeslot) error {
	return global.GVA_DB.Model(&model.VenueTimeslot{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *venueTimeslot) GetVenueTimeslot(id uint) (model.VenueTimeslot, error) {
	var res model.VenueTimeslot
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *venueTimeslot) GetVenueTimeslotList(req request.VenueTimeslotSearch) (list []model.VenueTimeslot, total int64, err error) {
	db := global.GVA_DB.Model(&model.VenueTimeslot{})
	if req.VenueID != nil {
		db = db.Where("venue_id = ?", *req.VenueID)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("venue_id ASC, start_time ASC").Find(&list).Error
	return
}

// GetVenueTimeslotsByVenue 获取某场地的全部时间段（用于下拉与前台）
func (s *venueTimeslot) GetVenueTimeslotsByVenue(venueID uint) (list []model.VenueTimeslot, err error) {
	err = global.GVA_DB.Where("venue_id = ?", venueID).Order("start_time ASC").Find(&list).Error
	return
}
