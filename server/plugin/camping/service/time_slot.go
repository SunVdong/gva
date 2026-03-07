package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
)

type timeSlot struct{}

func (s *timeSlot) CreateTimeSlot(m *model.CampingTimeSlot) error {
	return global.GVA_DB.Create(m).Error
}

func (s *timeSlot) DeleteTimeSlot(id uint) error {
	return global.GVA_DB.Delete(&model.CampingTimeSlot{}, "id = ?", id).Error
}

func (s *timeSlot) DeleteTimeSlotByIds(ids []uint) error {
	return global.GVA_DB.Delete(&[]model.CampingTimeSlot{}, "id in ?", ids).Error
}

func (s *timeSlot) UpdateTimeSlot(m model.CampingTimeSlot) error {
	return global.GVA_DB.Model(&model.CampingTimeSlot{}).Where("id = ?", m.ID).Updates(&m).Error
}

func (s *timeSlot) GetTimeSlot(id uint) (model.CampingTimeSlot, error) {
	var res model.CampingTimeSlot
	err := global.GVA_DB.Where("id = ?", id).First(&res).Error
	return res, err
}

func (s *timeSlot) GetTimeSlotList(req request.CampingTimeSlotSearch) (list []model.CampingTimeSlot, total int64, err error) {
	db := global.GVA_DB.Model(&model.CampingTimeSlot{})
	if err = db.Count(&total).Error; err != nil {
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("sort ASC, id ASC").Find(&list).Error
	return
}

// GetAllTimeSlots 获取全部时段（用于下拉与前台）
func (s *timeSlot) GetAllTimeSlots() (list []model.CampingTimeSlot, err error) {
	err = global.GVA_DB.Order("sort ASC, id ASC").Find(&list).Error
	return
}
