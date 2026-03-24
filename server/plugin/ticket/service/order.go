package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ticketOrder struct{}

func orderListToday() string {
	return time.Now().Format("2006-01-02")
}

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
	if req.OrderType != nil {
		today := orderListToday()
		switch *req.OrderType {
		case "pending_payment", "待支付":
			db = db.Where("status = ?", 0)
		case "pending_verify", "待核销":
			db = db.Where("status = ?", 1).
				Where("visit_date >= ?", today)
		case "completed", "已完成":
			db = db.Where("status IN (?)", []int{2, 3, 4, 5})
		}
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

// GetProductNamesByOrderIDs 批量获取订单对应的门票商品名称，返回 orderID -> productName
func (s *ticketOrder) GetProductNamesByOrderIDs(orderIDs []uint) (map[uint]string, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	type row struct {
		OrderID     uint   `gorm:"column:id"`
		ProductName string `gorm:"column:product_name"`
	}
	var rows []row
	if err := global.GVA_DB.Table(model.TicketOrder{}.TableName()).
		Select("orders.id, ticket_products.name as product_name").
		Joins("LEFT JOIN ticket_sku ON ticket_sku.id = orders.sku_id").
		Joins("LEFT JOIN ticket_products ON ticket_products.id = ticket_sku.product_id").
		Where("orders.id IN ?", orderIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		if strings.TrimSpace(r.ProductName) != "" {
			m[r.OrderID] = r.ProductName
		}
	}
	return m, nil
}

// GetScenicImageByOrderIDs 批量获取订单对应景区轮播图的第一张图片，返回 orderID -> imageURL
func (s *ticketOrder) GetScenicImageByOrderIDs(orderIDs []uint) (map[uint]string, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	type row struct {
		OrderID        uint   `gorm:"column:id"`
		CarouselImages string `gorm:"column:carousel_images"`
	}
	var rows []row
	if err := global.GVA_DB.Table(model.TicketOrder{}.TableName()).
		Select("orders.id, scenic_spots.carousel_images").
		Joins("LEFT JOIN ticket_sku ON ticket_sku.id = orders.sku_id").
		Joins("LEFT JOIN ticket_products ON ticket_products.id = ticket_sku.product_id").
		Joins("LEFT JOIN scenic_spots ON scenic_spots.id = ticket_products.scenic_id").
		Where("orders.id IN ?", orderIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		if _, ok := m[r.OrderID]; ok {
			continue
		}
		if r.CarouselImages == "" {
			continue
		}
		var imgs []string
		if err := json.Unmarshal([]byte(r.CarouselImages), &imgs); err == nil && len(imgs) > 0 {
			m[r.OrderID] = imgs[0]
		}
	}
	return m, nil
}

// OrderStatusLabel 根据订单游玩日计算展示状态
func (s *ticketOrder) OrderStatusLabel(order *model.TicketOrder) string {
	switch order.Status {
	case 0:
		return "待支付"
	case 1:
		today := orderListToday()
		if order.VisitDate.Format("2006-01-02") < today {
			return "已过期"
		}
		if order.VerifiedTimes > 0 {
			return "核销中"
		}
		return "待核销"
	case 2:
		return "已核销"
	case 3:
		return "已取消"
	case 4:
		return "已过期"
	case 5:
		return "已关闭"
	default:
		return "未知"
	}
}

// GetVerifyRecords 查询订单的核销记录列表
func (s *ticketOrder) GetVerifyRecords(orderID uint) (records []model.OrderVerifyRecord, err error) {
	err = global.GVA_DB.Where("order_id = ?", orderID).Order("verify_no ASC").Find(&records).Error
	return
}

func (s *ticketOrder) GetByID(id uint) (order model.TicketOrder, err error) {
	if err = global.GVA_DB.Where("id = ?", id).First(&order).Error; err != nil {
		return
	}
	s.fillProductName(&order)
	return
}

// fillProductName 补充订单对应的门票商品名称
func (s *ticketOrder) fillProductName(order *model.TicketOrder) {
	if order.SkuID == 0 {
		return
	}
	var sku model.TicketSku
	if e := global.GVA_DB.Where("id = ?", order.SkuID).First(&sku).Error; e != nil {
		return
	}
	var product model.TicketProduct
	if e := global.GVA_DB.Where("id = ?", sku.ProductID).First(&product).Error; e != nil {
		return
	}
	order.ProductName = product.Name
}

// GetByOrderNoPublic 根据订单号查询订单（公开给 H5 核销使用）
func (s *ticketOrder) GetByOrderNoPublic(orderNo string) (order model.TicketOrder, err error) {
	if err = global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return
	}
	s.fillProductName(&order)
	return
}

// validChineseMobile 中国大陆 11 位手机号：1 开头，第二位 3-9
var validChineseMobile = regexp.MustCompile(`^1[3-9]\d{9}$`)

// CreateOrder 小程序下单：生成订单号、校验 SKU 与库存、创建订单（userID 由 x-token 解析后传入）
func (s *ticketOrder) CreateOrder(userID uint, req request.MiniOrderCreate) (order model.TicketOrder, err error) {
	phone := strings.TrimSpace(req.BookerPhone)
	if !validChineseMobile.MatchString(phone) {
		return order, fmt.Errorf("预定人手机号格式不正确")
	}
	var orderNo string
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var sku model.TicketSku
		if e := tx.Where("id = ? AND status = ?", req.SkuID, 1).First(&sku).Error; e != nil {
			return fmt.Errorf("门票 SKU 不存在或已下架")
		}
		visitDate, e := time.ParseInLocation("2006-01-02", req.VisitDate, time.Local)
		if e != nil {
			return fmt.Errorf("游玩日期格式错误")
		}
		if sku.LimitBuy > 0 && req.Quantity > sku.LimitBuy {
			return fmt.Errorf("门票 %s 每单限购 %d 张", sku.Name, sku.LimitBuy)
		}
		var cal model.TicketCalendar
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("sku_id = ? AND visit_date = ? AND status = ?", req.SkuID, visitDate, 1).
			First(&cal).Error; e != nil {
			return fmt.Errorf("门票 %s 所选日期未开放或无库存", sku.Name)
		}
		ur := tx.Model(&model.TicketCalendar{}).
			Where("id = ? AND sold + ? <= stock", cal.ID, req.Quantity).
			UpdateColumn("sold", gorm.Expr("sold + ?", req.Quantity))
		if ur.Error != nil {
			return ur.Error
		}
		if ur.RowsAffected == 0 {
			return fmt.Errorf("门票 %s 所选日期库存不足", sku.Name)
		}
		totalAmount := sku.Price * float64(req.Quantity)
		useTimes := sku.UseTimes
		if useTimes <= 0 {
			useTimes = 1
		}
		orderNo = "T" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", rand.Intn(10000))
		order = model.TicketOrder{
			OrderNo:       orderNo,
			UserID:        userID,
			BookerName:    strings.TrimSpace(req.BookerName),
			BookerPhone:   phone,
			SkuID:         req.SkuID,
			SkuName:       sku.Name,
			Price:         sku.Price,
			Quantity:      req.Quantity,
			VisitDate:     visitDate,
			TotalAmount:   totalAmount,
			PayAmount:     totalAmount,
			Status:        0,
			TotalUseTimes: useTimes,
			VerifiedTimes: 0,
		}
		return tx.Create(&order).Error
	})
	if err != nil {
		return
	}
	err = global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error
	return
}

// VerifyOrder 核销订单（支持多次票累加核销，由后台或核销端调用）
func (s *ticketOrder) VerifyOrder(orderID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order model.TicketOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).First(&order).Error; err != nil || order.ID == 0 {
			return fmt.Errorf("订单不存在")
		}
		if order.Status == 2 {
			return fmt.Errorf("该订单已核销完毕")
		}
		if order.Status != 1 {
			return fmt.Errorf("仅待核销订单可核销")
		}
		totalUse := order.TotalUseTimes
		if totalUse <= 0 {
			totalUse = 1
		}
		if order.VerifiedTimes >= totalUse {
			return fmt.Errorf("该订单已核销完毕")
		}

		newVerified := order.VerifiedTimes + 1
		now := time.Now()

		record := model.OrderVerifyRecord{
			OrderID:    orderID,
			VerifyNo:   newVerified,
			VerifiedAt: now,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		updates := map[string]any{
			"verified_times": newVerified,
		}
		if newVerified >= totalUse {
			updates["verified_at"] = now
			updates["status"] = 2
		}
		return tx.Model(&model.TicketOrder{}).Where("id = ?", orderID).Updates(updates).Error
	})
}

// VerifyOrderByOrderNoPublic 根据订单号核销订单（公开给 H5 核销使用）
func (s *ticketOrder) VerifyOrderByOrderNoPublic(orderNo string) error {
	var order model.TicketOrder
	if err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil || order.ID == 0 {
		return fmt.Errorf("订单不存在")
	}
	return s.VerifyOrder(order.ID)
}
