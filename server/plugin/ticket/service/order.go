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
	"gorm.io/gorm/clause"
)

type ticketOrder struct{}

// orderListToday 用于与订单项最大游玩日比较，取当地日期
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
	// 订单类型：待支付、待核销、已完成（不传则不过滤）
	if req.OrderType != nil {
		today := orderListToday()
		switch *req.OrderType {
		case "pending_payment", "待支付":
			// 待支付：状态为0（待支付）
			db = db.Where("status = ?", 0)
		case "pending_verify", "待核销":
			// 待核销：状态为1（待核销），且最晚游玩日 >= 今天（尚未过期）
			db = db.Where("status = ?", 1).
				Where("(SELECT COALESCE(MAX(visit_date),'1970-01-01') FROM order_items WHERE order_items.order_id = orders.id) >= ?", today)
		case "completed", "已完成":
			// 已完成：已核销、已取消、已过期、已关闭
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

// GetMaxVisitDateByOrderIDs 批量获取订单项中最晚游玩日期，用于计算订单状态文案。返回 orderID -> "2006-01-02"
func (s *ticketOrder) GetMaxVisitDateByOrderIDs(orderIDs []uint) (map[uint]string, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	type row struct {
		OrderID  uint      `gorm:"column:order_id"`
		MaxVisit time.Time `gorm:"column:max_visit"`
	}
	var rows []row
	err := global.GVA_DB.Model(&model.OrderItem{}).Select("order_id, MAX(visit_date) as max_visit").
		Where("order_id IN ?", orderIDs).Group("order_id").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		m[r.OrderID] = r.MaxVisit.Format("2006-01-02")
	}
	return m, nil
}

// OrderStatusLabel 根据订单与最晚游玩日计算展示状态：待支付、待核销、已核销、已取消、已过期、已关闭
func (s *ticketOrder) OrderStatusLabel(order *model.TicketOrder, maxVisitDate string) string {
	switch order.Status {
	case 0:
		return "待支付"
	case 1:
		// 待核销：根据最晚游玩日判断是否已过期
		today := orderListToday()
		if maxVisitDate != "" && maxVisitDate < today {
			return "已过期"
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

func (s *ticketOrder) GetByID(id uint) (order model.TicketOrder, items []model.OrderItem, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&order).Error
	if err != nil {
		return
	}
	// 查询订单项
	if err = global.GVA_DB.Where("order_id = ?", id).Order("id").Find(&items).Error; err != nil {
		return
	}
	if err = s.fillOrderItemsProductName(items); err != nil {
		return
	}
	return
}

// fillOrderItemsProductName 补充每个订单项对应的门票商品名称（ticket_products.name）
func (s *ticketOrder) fillOrderItemsProductName(items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	// 收集 SKU ID
	skuIDs := make([]uint, 0, len(items))
	for _, it := range items {
		skuIDs = append(skuIDs, it.SkuID)
	}

	// 查询 SKU，拿到 ProductID
	var skus []model.TicketSku
	if e := global.GVA_DB.Where("id IN ?", skuIDs).Find(&skus).Error; e != nil || len(skus) == 0 {
		return e
	}
	productIDs := make([]uint, 0, len(skus))
	for _, sku := range skus {
		productIDs = append(productIDs, sku.ProductID)
	}

	// 查询门票商品，拿到商品名称
	var products []model.TicketProduct
	if e := global.GVA_DB.Where("id IN ?", productIDs).Find(&products).Error; e != nil {
		return e
	}
	if len(products) == 0 {
		return nil
	}
	productNameByID := make(map[uint]string, len(products))
	for _, p := range products {
		productNameByID[p.ID] = p.Name
	}

	// 构建 SKU -> 商品名称 的映射
	productNameBySkuID := make(map[uint]string, len(skus))
	for _, sku := range skus {
		if name, ok := productNameByID[sku.ProductID]; ok {
			productNameBySkuID[sku.ID] = name
		}
	}

	// 写回每个订单项的 ProductName 字段
	for i := range items {
		if name, ok := productNameBySkuID[items[i].SkuID]; ok {
			items[i].ProductName = name
		}
	}
	return nil
}

// GetByOrderNoPublic 根据订单号查询订单及订单项（公开给 H5 核销使用）
func (s *ticketOrder) GetByOrderNoPublic(orderNo string) (order model.TicketOrder, items []model.OrderItem, err error) {
	if err = global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return
	}
	if err = global.GVA_DB.Where("order_id = ?", order.ID).Order("id").Find(&items).Error; err != nil {
		return
	}
	if err = s.fillOrderItemsProductName(items); err != nil {
		return
	}
	return
}

// validChineseMobile 中国大陆 11 位手机号：1 开头，第二位 3-9
var validChineseMobile = regexp.MustCompile(`^1[3-9]\d{9}$`)

// CreateOrder 小程序下单：生成订单号、校验 SKU 与库存、创建订单及订单项（userID 由 x-token 解析后传入）
func (s *ticketOrder) CreateOrder(userID uint, req request.MiniOrderCreate) (order model.TicketOrder, err error) {
	phone := strings.TrimSpace(req.BookerPhone)
	if !validChineseMobile.MatchString(phone) {
		return order, fmt.Errorf("预定人手机号格式不正确")
	}
	var orderNo string
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []model.OrderItem
		// 每单限购：同一订单内，同一 SKU 的数量汇总后，不得超过 sku.limit_buy（0 表示不限购）
		perOrderQty := make(map[uint]int)
		for _, it := range req.Items {
			var sku model.TicketSku
			if e := tx.Where("id = ? AND status = ?", it.SkuID, 1).First(&sku).Error; e != nil {
				return fmt.Errorf("门票 SKU 不存在或已下架")
			}
			visitDate, e := time.ParseInLocation("2006-01-02", it.VisitDate, time.Local)
			if e != nil {
				return fmt.Errorf("游玩日期格式错误")
			}
			perOrderQty[it.SkuID] += it.Quantity
			if sku.LimitBuy > 0 && perOrderQty[it.SkuID] > sku.LimitBuy {
				return fmt.Errorf("门票 %s 每单限购 %d 张", sku.Name, sku.LimitBuy)
			}
			// 日历库存校验：必须存在当日排期才可售，且按 stock-sold 判断可售数量
			var cal model.TicketCalendar
			if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("sku_id = ? AND visit_date = ? AND status = ?", it.SkuID, visitDate, 1).
				First(&cal).Error; e != nil {
				return fmt.Errorf("门票 %s 所选日期未开放或无库存", sku.Name)
			}
			// 扣减库存：sold += quantity（带条件避免并发超卖）
			ur := tx.Model(&model.TicketCalendar{}).
				Where("id = ? AND sold + ? <= stock", cal.ID, it.Quantity).
				UpdateColumn("sold", gorm.Expr("sold + ?", it.Quantity))
			if ur.Error != nil {
				return ur.Error
			}
			if ur.RowsAffected == 0 {
				return fmt.Errorf("门票 %s 所选日期库存不足", sku.Name)
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
			UserID:      userID,
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

// VerifyOrder 核销订单（仅待核销状态的订单可核销，由后台或核销端调用）
func (s *ticketOrder) VerifyOrder(orderID uint) error {
	var order model.TicketOrder
	if err := global.GVA_DB.Where("id = ?", orderID).First(&order).Error; err != nil || order.ID == 0 {
		return fmt.Errorf("订单不存在")
	}
	// 兼容历史数据：可能 status 已更新但 verified_at 为空，或反之
	if order.Status == 2 || order.VerifiedAt != nil {
		return fmt.Errorf("该订单已核销")
	}
	if order.Status != 1 {
		return fmt.Errorf("仅待核销订单可核销")
	}
	now := time.Now()
	return global.GVA_DB.Model(&model.TicketOrder{}).
		Where("id = ?", orderID).
		Updates(map[string]any{
			"verified_at": now,
			"status":      2,
		}).Error
}

// VerifyOrderByOrderNoPublic 根据订单号核销订单（公开给 H5 核销使用）
func (s *ticketOrder) VerifyOrderByOrderNoPublic(orderNo string) error {
	var order model.TicketOrder
	if err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil || order.ID == 0 {
		return fmt.Errorf("订单不存在")
	}
	return s.VerifyOrder(order.ID)
}
