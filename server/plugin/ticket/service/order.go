package service

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"gorm.io/gorm"
)

type ticketOrder struct{}

func (s *ticketOrder) GetList(req request.TicketOrderSearch) (list []model.TicketOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketOrder{})
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
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
	err = db.Order("id DESC").Find(&list).Error
	return
}

func (s *ticketOrder) GetByID(id uint) (order model.TicketOrder, items []model.OrderItem, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&order).Error
	if err != nil {
		return
	}
	err = global.GVA_DB.Where("order_id = ?", id).Order("id").Find(&items).Error
	return
}

// validChineseMobile 中国大陆 11 位手机号：1 开头，第二位 3-9
var validChineseMobile = regexp.MustCompile(`^1[3-9]\d{9}$`)

// CreateOrder 小程序下单：生成订单号、校验 SKU 与库存、创建订单及订单项
func (s *ticketOrder) CreateOrder(req request.MiniOrderCreate) (order model.TicketOrder, err error) {
	phone := strings.TrimSpace(req.BookerPhone)
	if !validChineseMobile.MatchString(phone) {
		return order, fmt.Errorf("预定人手机号格式不正确")
	}
	var orderNo string
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []model.OrderItem
		for _, it := range req.Items {
			var sku model.TicketSku
			if e := tx.Where("id = ? AND status = ?", it.SkuID, 1).First(&sku).Error; e != nil {
				return fmt.Errorf("门票 SKU 不存在或已下架")
			}
			visitDate, e := time.ParseInLocation("2006-01-02", it.VisitDate, time.Local)
			if e != nil {
				return fmt.Errorf("游玩日期格式错误")
			}
			var cal model.TicketCalendar
			if e := tx.Where("sku_id = ? AND visit_date = ? AND status = ?", it.SkuID, visitDate, 1).First(&cal).Error; e == nil {
				if cal.Stock-cal.Sold < it.Quantity {
					return fmt.Errorf("门票 %s 所选日期库存不足", sku.Name)
				}
			}
			subTotal := sku.Price * float64(it.Quantity)
			totalAmount += subTotal
			orderItems = append(orderItems, model.OrderItem{
				SkuID:     it.SkuID,
				SkuName:   sku.Name,
				Price:     sku.Price,
				Quantity:  it.Quantity,
				VisitDate: visitDate,
			})
		}
		orderNo = "T" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", rand.Intn(10000))
		order = model.TicketOrder{
			OrderNo:     orderNo,
			UserID:      req.UserID,
			BookerName:  strings.TrimSpace(req.BookerName),
			BookerPhone: phone,
			TotalAmount: totalAmount,
			PayAmount:   totalAmount,
			Status:      0,
		}
		if e := tx.Create(&order).Error; e != nil {
			return e
		}
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
			if e := tx.Create(&orderItems[i]).Error; e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	err = global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error
	return
}
