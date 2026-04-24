package service

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
			// 使用本地时区解析，得到 "YYYY-MM-DD 00:00:00 Local"。
			// 配合驱动 loc=Local，到达 MySQL 时仍为本地零点，DATE 列能精确匹配，不会产生小时偏移。
			t, err := time.ParseInLocation("2006-01-02", item.VisitDate, time.Local)
			if err != nil {
				return fmt.Errorf("无效的日期格式: %s", item.VisitDate)
			}

			// 使用 INSERT ... ON DUPLICATE KEY UPDATE：
			//   1) 原子操作，避免并发双写产生 Duplicate entry；
			//   2) MySQL 唯一索引不感知软删除，冲突即走 UPDATE，并把 deleted_at 置 NULL 以"复活"旧记录；
			//   3) DoUpdates 不包含 sold，避免覆盖已售数量。
			err = tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "sku_id"},
					{Name: "visit_date"},
				},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"stock":      item.Stock,
					"status":     item.Status,
					"deleted_at": nil,
					"updated_at": time.Now(),
				}),
			}).Create(&model.TicketCalendar{
				SkuID:     item.SkuID,
				VisitDate: t,
				Stock:     item.Stock,
				Sold:      0,
				Status:    item.Status,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
