package service

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type reservation struct{}

func (s *reservation) CreateReservation(req request.CreateVenueReservationRequest, userID uint) (model.VenueReservation, error) {
	reserveDate, err := time.ParseInLocation("2006-01-02", req.ReserveDate, time.Local)
	if err != nil {
		return model.VenueReservation{}, fmt.Errorf("预约日期格式错误")
	}
	// 检查日历该日是否可约
	open, err := Service.VenueCalendar.IsDateOpen(req.VenueID, reserveDate)
	if err != nil || !open {
		return model.VenueReservation{}, fmt.Errorf("该日期暂不开放预约")
	}
	// 检查时间段是否存在且属于该场地
	var slot model.VenueTimeslot
	if err := global.GVA_DB.Where("id = ? AND venue_id = ?", req.TimeslotID, req.VenueID).First(&slot).Error; err != nil {
		return model.VenueReservation{}, fmt.Errorf("时间段无效")
	}
	// 检查已预约数是否已达 capacity
	var count int64
	global.GVA_DB.Model(&model.VenueReservation{}).
		Where("venue_id = ? AND reserve_date = ? AND timeslot_id = ? AND status IN (0, 1)", req.VenueID, reserveDate, req.TimeslotID).
		Count(&count)
	if int(count) >= slot.Capacity {
		return model.VenueReservation{}, fmt.Errorf("该时段已约满")
	}
	reservationNo := fmt.Sprintf("R%d%s", time.Now().UnixNano()/1e6, utils.RandomString(6))
	code := utils.RandomString(12)
	for {
		var c int64
		global.GVA_DB.Model(&model.VenueReservation{}).Where("verify_code = ?", code).Count(&c)
		if c == 0 {
			break
		}
		code = utils.RandomString(12)
	}
	m := model.VenueReservation{
		ReservationNo: reservationNo,
		UserID:        userID,
		VenueID:       req.VenueID,
		TimeslotID:    req.TimeslotID,
		ReserveDate:   reserveDate,
		ContactName:   req.ContactName,
		ContactPhone:  req.ContactPhone,
		ContactCount:  req.ContactCount,
		Status:        0,
		VerifyCode:    code,
	}
	return m, global.GVA_DB.Create(&m).Error
}

func (s *reservation) GetReservation(id uint) (model.VenueReservation, error) {
	var res model.VenueReservation
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *reservation) GetReservationByVerifyCode(code string) (model.VenueReservation, error) {
	var res model.VenueReservation
	err := global.GVA_DB.Where("verify_code = ?", code).First(&res).Error
	return res, err
}

func (s *reservation) GetReservationList(req request.VenueReservationSearch) (list []model.VenueReservation, total int64, err error) {
	db := global.GVA_DB.Model(&model.VenueReservation{})
	if req.VenueID != nil {
		db = db.Where("venue_id = ?", *req.VenueID)
	}
	if req.ReserveDate != nil {
		db = db.Where("reserve_date = ?", *req.ReserveDate)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.VerifyCode != "" {
		db = db.Where("verify_code = ?", req.VerifyCode)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("reserve_date DESC, timeslot_id ASC").Find(&list).Error
	return
}

func (s *reservation) VerifyReservation(id uint) error {
	return global.GVA_DB.Model(&model.VenueReservation{}).Where("id = ? AND status = 0", id).Update("status", 1).Error
}

func (s *reservation) VerifyReservationByCode(code string) error {
	return global.GVA_DB.Model(&model.VenueReservation{}).Where("verify_code = ? AND status = 0", code).Update("status", 1).Error
}

func (s *reservation) CancelReservation(id uint) error {
	return global.GVA_DB.Model(&model.VenueReservation{}).Where("id = ?", id).Update("status", 2).Error
}

// GetReservedTimeslotIds 获取某场地某日已约满的时段ID列表（已预约数 >= capacity 的时段）
func (s *reservation) GetReservedTimeslotIds(venueID uint, reserveDate time.Time) ([]uint, error) {
	var slots []model.VenueTimeslot
	if err := global.GVA_DB.Where("venue_id = ?", venueID).Find(&slots).Error; err != nil {
		return nil, err
	}
	var fullIds []uint
	for _, slot := range slots {
		var count int64
		global.GVA_DB.Model(&model.VenueReservation{}).
			Where("venue_id = ? AND reserve_date = ? AND timeslot_id = ? AND status IN (0, 1)", venueID, reserveDate, slot.ID).
			Count(&count)
		if int(count) >= slot.Capacity {
			fullIds = append(fullIds, slot.ID)
		}
	}
	return fullIds, nil
}
