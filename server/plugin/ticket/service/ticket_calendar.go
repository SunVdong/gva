package service

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"gorm.io/gorm"
)

type ticketCalendar struct{}

func (s *ticketCalendar) GetBySku(req request.TicketCalendarSearch) (list []model.TicketCalendar, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketCalendar{}).Where("sku_id = ?", req.SkuID)
	if req.VisitDate != "" {
		db = db.Where("visit_date = ?", req.VisitDate)
	}
	if req.StartDate != "" && req.EndDate != "" {
		db = db.Where("visit_date >= ? AND visit_date <= ?", req.StartDate, req.EndDate)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 31
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := pageSize * (page - 1)
	db = db.Limit(pageSize).Offset(offset).Order("visit_date ASC")
	err = db.Find(&list).Error
	return
}

func (s *ticketCalendar) Set(req request.TicketCalendarSet) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.List {
			t, _ := time.Parse("2006-01-02", item.VisitDate)
			var cal model.TicketCalendar
			err := tx.Where("sku_id = ? AND visit_date = ?", item.SkuID, t).First(&cal).Error
			if err == nil {
				cal.Stock = item.Stock
				cal.Status = item.Status
				if err := tx.Save(&cal).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&model.TicketCalendar{
					SkuID:     item.SkuID,
					VisitDate: t,
					Stock:     item.Stock,
					Sold:      0,
					Status:    item.Status,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
