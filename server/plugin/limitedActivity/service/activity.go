package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
)

type activity struct{}

func (s *activity) Create(m *model.Activity) error {
	m.Name = strings.TrimSpace(m.Name)
	m.Address = strings.TrimSpace(m.Address)
	if m.Name == "" {
		return fmt.Errorf("活动名称不能为空")
	}
	if m.Quota < 0 {
		return fmt.Errorf("名额不能为负数")
	}
	if m.Price < 0 || m.MarketPrice < 0 {
		return fmt.Errorf("价格不能为负数")
	}
	if !m.EndTime.After(m.StartTime) {
		return fmt.Errorf("结束时间必须晚于开始时间")
	}
	m.Sold = 0
	return global.GVA_DB.Create(m).Error
}

func (s *activity) Delete(id uint) error {
	var cnt int64
	if err := global.GVA_DB.Model(&model.ActivityOrder{}).
		Where("activity_id = ? AND status IN ?", id, []int{0, 1, 7}).
		Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return fmt.Errorf("存在待支付/待核销/退款中订单，无法删除活动")
	}
	return global.GVA_DB.Delete(&model.Activity{}, "id = ?", id).Error
}

func (s *activity) DeleteByIds(ids []uint) error {
	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

func (s *activity) Update(m model.Activity) error {
	var old model.Activity
	if err := global.GVA_DB.Where("id = ?", m.ID).First(&old).Error; err != nil {
		return fmt.Errorf("活动不存在")
	}
	m.Name = strings.TrimSpace(m.Name)
	m.Address = strings.TrimSpace(m.Address)
	if m.Name == "" {
		return fmt.Errorf("活动名称不能为空")
	}
	if m.Quota < old.Sold {
		return fmt.Errorf("名额不能小于已占用人数(%d)", old.Sold)
	}
	if m.Price < 0 || m.MarketPrice < 0 {
		return fmt.Errorf("价格不能为负数")
	}
	if !m.EndTime.After(m.StartTime) {
		return fmt.Errorf("结束时间必须晚于开始时间")
	}
	// 不允许通过更新接口直接改 sold
	updates := map[string]any{
		"name":         m.Name,
		"address":      m.Address,
		"detail":       m.Detail,
		"start_time":   m.StartTime,
		"end_time":     m.EndTime,
		"market_price": m.MarketPrice,
		"price":        m.Price,
		"quota":        m.Quota,
		"cover_image":  m.CoverImage,
		"group_qr":     m.GroupQr,
		"service_qr":   m.ServiceQr,
		"sort":         m.Sort,
		"status":       m.Status,
		"updated_by":   m.UpdatedBy,
	}
	return global.GVA_DB.Model(&model.Activity{}).Where("id = ?", m.ID).Updates(updates).Error
}

func (s *activity) Get(id uint) (model.Activity, error) {
	var res model.Activity
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	if err == nil {
		s.fillVirtual(&res)
	}
	return res, err
}

func (s *activity) GetList(req request.ActivitySearch) (list []model.Activity, total int64, err error) {
	db := global.GVA_DB.Model(&model.Activity{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("sort ASC, id ASC").Find(&list).Error
	for i := range list {
		s.fillVirtual(&list[i])
	}
	return
}

// GetMiniList 小程序可见活动列表（仅显示中，按 sort ASC, id ASC）
func (s *activity) GetMiniList(req request.ActivitySearch) (list []model.Activity, total int64, err error) {
	db := global.GVA_DB.Model(&model.Activity{}).Where("status = ?", 1)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("sort ASC, id ASC").Find(&list).Error
	for i := range list {
		s.fillVirtual(&list[i])
	}
	return
}

// GetMiniDetail 小程序活动详情（仅显示中）
func (s *activity) GetMiniDetail(id uint) (model.Activity, error) {
	var res model.Activity
	err := global.GVA_DB.Where("id = ? AND status = ?", id, 1).First(&res).Error
	if err == nil {
		s.fillVirtual(&res)
	}
	return res, err
}

func (s *activity) fillVirtual(a *model.Activity) {
	remaining := a.Quota - a.Sold
	if remaining < 0 {
		remaining = 0
	}
	a.Remaining = remaining
	now := time.Now()
	a.CanSignup = a.Status == 1 && !now.Before(a.StartTime) && !now.After(a.EndTime) && remaining > 0
}
