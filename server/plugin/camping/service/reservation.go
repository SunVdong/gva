package service

import (
	"fmt"
	"regexp"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

// 中国大陆手机号：1 开头，第二位 3-9，共 11 位数字
var reChineseMobile = regexp.MustCompile(`^1[3-9]\d{9}$`)

type reservation struct{}

func (s *reservation) CreateReservation(req request.CreateVenueReservationRequest, userID uint) (model.VenueReservation, error) {
	if !reChineseMobile.MatchString(req.ContactPhone) {
		return model.VenueReservation{}, fmt.Errorf("请输入正确的手机号")
	}
	reserveDate, err := time.ParseInLocation("2006-01-02", req.ReserveDate, time.Local)
	if err != nil {
		return model.VenueReservation{}, fmt.Errorf("预约日期格式错误")
	}
	reserveDate = dateOnly(reserveDate)
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

// UpdateReservation 修改预约
// 仅允许修改未取消的预约，且必须属于当前用户；会重新检查日期开放情况、时段有效性与容量。
func (s *reservation) UpdateReservation(req request.UpdateVenueReservationRequest, userID uint) (model.VenueReservation, error) {
	if !reChineseMobile.MatchString(req.ContactPhone) {
		return model.VenueReservation{}, fmt.Errorf("请输入正确的手机号")
	}
	// 查询原预约
	var original model.VenueReservation
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&original).Error; err != nil {
		return model.VenueReservation{}, fmt.Errorf("预约不存在")
	}
	if original.UserID != userID {
		return model.VenueReservation{}, fmt.Errorf("无权修改该预约")
	}
	if original.Status == 2 {
		return model.VenueReservation{}, fmt.Errorf("该预约已取消，无法修改")
	}

	// 解析并规范化预约日期
	reserveDate, err := time.ParseInLocation("2006-01-02", req.ReserveDate, time.Local)
	if err != nil {
		return model.VenueReservation{}, fmt.Errorf("预约日期格式错误")
	}
	reserveDate = dateOnly(reserveDate)

	// 检查日历该日是否可约
	open, err := Service.VenueCalendar.IsDateOpen(original.VenueID, reserveDate)
	if err != nil || !open {
		return model.VenueReservation{}, fmt.Errorf("该日期暂不开放预约")
	}

	// 检查时间段是否存在且属于该场地
	var slot model.VenueTimeslot
	if err := global.GVA_DB.Where("id = ? AND venue_id = ?", req.TimeslotID, original.VenueID).First(&slot).Error; err != nil {
		return model.VenueReservation{}, fmt.Errorf("时间段无效")
	}

	// 检查已预约数是否已达 capacity（排除当前这条预约）
	var count int64
	global.GVA_DB.Model(&model.VenueReservation{}).
		Where("venue_id = ? AND reserve_date = ? AND timeslot_id = ? AND status IN (0, 1) AND id <> ?", original.VenueID, reserveDate, req.TimeslotID, original.ID).
		Count(&count)
	if int(count) >= slot.Capacity {
		return model.VenueReservation{}, fmt.Errorf("该时段已约满")
	}

	// 更新字段
	original.ReserveDate = reserveDate
	original.TimeslotID = req.TimeslotID
	original.ContactName = req.ContactName
	original.ContactPhone = req.ContactPhone
	original.ContactCount = req.ContactCount

	return original, global.GVA_DB.Save(&original).Error
}

func (s *reservation) GetReservation(id uint) (model.VenueReservation, error) {
	var res model.VenueReservation
	if err := global.GVA_DB.Where("id = ?", id).First(&res).Error; err != nil {
		return res, err
	}
	_ = s.maybeExpireReservation(&res)
	return res, nil
}

func (s *reservation) GetReservationByVerifyCode(code string) (model.VenueReservation, error) {
	var res model.VenueReservation
	err := global.GVA_DB.Where("verify_code = ?", code).First(&res).Error
	return res, err
}

func (s *reservation) GetReservationList(req request.VenueReservationSearch) (list []model.VenueReservation, total int64, err error) {
	db := global.GVA_DB.Model(&model.VenueReservation{})
	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.VenueID != nil {
		db = db.Where("venue_id = ?", *req.VenueID)
	}
	if req.ReserveDate != nil {
		d := dateOnly(*req.ReserveDate)
		db = db.Where("reserve_date = ?", d)
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
	if err = db.Order("reserve_date DESC, timeslot_id ASC").Find(&list).Error; err != nil {
		return
	}
	for i := range list {
		_ = s.maybeExpireReservation(&list[i])
	}
	return
}

func (s *reservation) VerifyReservation(id uint) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.VenueReservation{}).
		Where("id = ? AND status = 0", id).
		Updates(map[string]interface{}{"status": 1, "verified_at": now}).Error
}

func (s *reservation) VerifyReservationByCode(code string) error {
	now := time.Now()
	res := global.GVA_DB.Model(&model.VenueReservation{}).
		Where("verify_code = ? AND status = 0", code).
		Updates(map[string]interface{}{"status": 1, "verified_at": now})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("核销码无效或该预约已核销/已取消/已过期")
	}
	return nil
}

func (s *reservation) CancelReservation(id uint) error {
	var res model.VenueReservation
	if err := global.GVA_DB.Where("id = ?", id).First(&res).Error; err != nil {
		return err
	}
	// 校验：距开始时间不足 N 小时不可取消（N 为场地的 RefundChangeHours）
	var venue model.Venue
	if err := global.GVA_DB.Where("id = ?", res.VenueID).First(&venue).Error; err != nil {
		return err
	}
	if venue.RefundChangeHours > 0 {
		var slot model.VenueTimeslot
		if err := global.GVA_DB.Where("id = ?", res.TimeslotID).First(&slot).Error; err != nil {
			return err
		}
		startTime, err := combineDateAndTimeOnly(res.ReserveDate, slot.StartTime)
		if err != nil {
			return err
		}
		now := time.Now()
		minCancelTime := startTime.Add(-time.Duration(venue.RefundChangeHours) * time.Hour)
		if now.After(minCancelTime) {
			return fmt.Errorf("距开始时间不足 %d 小时，不可取消", venue.RefundChangeHours)
		}
	}
	return global.GVA_DB.Model(&model.VenueReservation{}).Where("id = ?", id).Update("status", 2).Error
}

// maybeExpireReservation 如果预约为待核销且当前时间已超过该场次结束时间，则将其标记为已过期（status=3）
func (s *reservation) maybeExpireReservation(res *model.VenueReservation) error {
	if res == nil {
		return nil
	}
	// 仅对待核销状态进行过期判断
	if res.Status != 0 {
		return nil
	}
	// 查询时段，优先使用结束时间，若为空则使用开始时间
	var slot model.VenueTimeslot
	if err := global.GVA_DB.Where("id = ?", res.TimeslotID).First(&slot).Error; err != nil {
		return nil
	}
	var endBase model.TimeOnly
	if string(slot.EndTime) != "" {
		endBase = slot.EndTime
	} else {
		endBase = slot.StartTime
	}
	// 如果仍然为空，则无法判断
	if string(endBase) == "" {
		return nil
	}
	endTime, err := combineDateAndTimeOnly(res.ReserveDate, endBase)
	if err != nil {
		return nil
	}
	if time.Now().After(endTime) {
		// 更新数据库状态为已过期
		if e := global.GVA_DB.Model(&model.VenueReservation{}).
			Where("id = ? AND status = 0", res.ID).
			Update("status", 3).Error; e != nil {
			return e
		}
		res.Status = 3
	}
	return nil
}

// combineDateAndTimeOnly 将日期与 TimeOnly 拼成当地时区的 time.Time（使用 time.Local 便于与当前时间比较）
func combineDateAndTimeOnly(date time.Time, t model.TimeOnly) (time.Time, error) {
	s := string(t)
	loc := time.Local
	if s == "" {
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), nil
	}
	var tParsed time.Time
	var err error
	if len(s) >= 8 {
		tParsed, err = time.ParseInLocation("15:04:05", s, loc)
	} else {
		tParsed, err = time.ParseInLocation("15:04", s, loc)
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), tParsed.Hour(), tParsed.Minute(), tParsed.Second(), 0, loc), nil
}

// GetReservedTimeslotIds 获取某场地某日已约满的时段ID列表（已预约数 >= capacity 的时段）
func (s *reservation) GetReservedTimeslotIds(venueID uint, reserveDate time.Time) ([]uint, error) {
	d := dateOnly(reserveDate)
	var slots []model.VenueTimeslot
	if err := global.GVA_DB.Where("venue_id = ?", venueID).Find(&slots).Error; err != nil {
		return nil, err
	}
	var fullIds []uint
	for _, slot := range slots {
		var count int64
		global.GVA_DB.Model(&model.VenueReservation{}).
			Where("venue_id = ? AND reserve_date = ? AND timeslot_id = ? AND status IN (0, 1)", venueID, d, slot.ID).
			Count(&count)
		if int(count) >= slot.Capacity {
			fullIds = append(fullIds, slot.ID)
		}
	}
	return fullIds, nil
}
