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

func (s *reservation) CreateReservation(req request.CreateCampingReservationRequest) (model.CampingReservation, error) {
	reserveDate, err := time.ParseInLocation("2006-01-02", req.ReserveDate, time.Local)
	if err != nil {
		return model.CampingReservation{}, fmt.Errorf("预约日期格式错误")
	}
	// 检查同一场地同一日期相同时段是否已被预约
	var count int64
	global.GVA_DB.Model(&model.CampingReservation{}).Where("site_id = ? AND reserve_date = ? AND time_slot_id = ? AND status IN (0,1)", req.SiteID, reserveDate, req.TimeSlotID).Count(&count)
	if count > 0 {
		return model.CampingReservation{}, fmt.Errorf("该时段已被预约")
	}
	code := utils.RandomString(12)
	for {
		var c int64
		global.GVA_DB.Model(&model.CampingReservation{}).Where("verify_code = ?", code).Count(&c)
		if c == 0 {
			break
		}
		code = utils.RandomString(12)
	}
	m := model.CampingReservation{
		SiteID:      req.SiteID,
		ReserveDate: reserveDate,
		TimeSlotID:  req.TimeSlotID,
		BookerName:  req.BookerName,
		Phone:       req.Phone,
		PeopleCount: req.PeopleCount,
		Remark:      req.Remark,
		VerifyCode:  code,
		Status:      0,
	}
	return m, global.GVA_DB.Create(&m).Error
}

func (s *reservation) GetReservation(id uint) (model.CampingReservation, error) {
	var res model.CampingReservation
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *reservation) GetReservationByVerifyCode(code string) (model.CampingReservation, error) {
	var res model.CampingReservation
	err := global.GVA_DB.Where("verify_code = ?", code).First(&res).Error
	return res, err
}

func (s *reservation) GetReservationList(req request.CampingReservationSearch) (list []model.CampingReservation, total int64, err error) {
	db := global.GVA_DB.Model(&model.CampingReservation{})
	if req.SiteID != nil {
		db = db.Where("site_id = ?", *req.SiteID)
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
	err = db.Order("reserve_date DESC, time_slot_id ASC").Find(&list).Error
	return
}

func (s *reservation) VerifyReservation(id uint) error {
	return global.GVA_DB.Model(&model.CampingReservation{}).Where("id = ? AND status = 0", id).Update("status", 1).Error
}

func (s *reservation) VerifyReservationByCode(code string) error {
	return global.GVA_DB.Model(&model.CampingReservation{}).Where("verify_code = ? AND status = 0", code).Update("status", 1).Error
}

func (s *reservation) CancelReservation(id uint) error {
	return global.GVA_DB.Model(&model.CampingReservation{}).Where("id = ?", id).Update("status", 2).Error
}

// GetReservedSlotIds 获取某场地某日已被预约的时段ID列表
func (s *reservation) GetReservedSlotIds(siteID uint, reserveDate time.Time) ([]uint, error) {
	var ids []uint
	err := global.GVA_DB.Model(&model.CampingReservation{}).Where("site_id = ? AND reserve_date = ? AND status IN (0, 1)", siteID, reserveDate).Pluck("time_slot_id", &ids).Error
	return ids, err
}
