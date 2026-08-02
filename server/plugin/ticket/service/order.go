package service

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	miniPay "github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ticketOrder struct{}

func orderListToday() string {
	return time.Now().Format("2006-01-02")
}

// orderPaidUseTimes 付费次数 A；老单缺快照时按 totalUseTimes 兼容
func orderPaidUseTimes(order *model.TicketOrder) int {
	if order.PaidUseTimes > 0 {
		return order.PaidUseTimes
	}
	if order.TotalUseTimes > 0 {
		return order.TotalUseTimes
	}
	return 1
}

// RemainingPaidUseTimes 剩余付费可核销次数
func (s *ticketOrder) RemainingPaidUseTimes(order *model.TicketOrder) int {
	a := orderPaidUseTimes(order)
	paidConsumed := order.VerifiedTimes
	if paidConsumed > a {
		paidConsumed = a
	}
	remaining := a - paidConsumed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// CalcRefundFen 按剩余付费次数比例计算退款分：Round(payAmount*100*remainingPaid/A)
func (s *ticketOrder) CalcRefundFen(order *model.TicketOrder) (refundFen, totalFen, remainingPaid int, err error) {
	a := orderPaidUseTimes(order)
	if a <= 0 {
		return 0, 0, 0, fmt.Errorf("订单付费次数异常")
	}
	remainingPaid = s.RemainingPaidUseTimes(order)
	if remainingPaid <= 0 {
		return 0, 0, 0, fmt.Errorf("付费次数已用尽，不可退款")
	}
	totalFen = int(math.Round(order.PayAmount * 100))
	if totalFen <= 0 {
		return 0, 0, 0, fmt.Errorf("订单金额异常")
	}
	refundFen = int(math.Round(order.PayAmount * 100 * float64(remainingPaid) / float64(a)))
	if refundFen <= 0 {
		return 0, 0, 0, fmt.Errorf("退款金额异常")
	}
	if refundFen > totalFen {
		refundFen = totalFen
	}
	return
}

// AdminRefund 后台按付费次数剩余比例发起微信退款
func (s *ticketOrder) AdminRefund(orderID uint) error {
	return s.RequestRefund(orderID, 0)
}

// RequestRefund 发起退款。userID>0 时校验本人订单；先落库 status=7 再调微信。
func (s *ticketOrder) RequestRefund(orderID uint, userID uint) error {
	var order model.TicketOrder
	var refundFen, totalFen int
	var refundNo string

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).First(&order).Error; e != nil || order.ID == 0 {
			return fmt.Errorf("订单不存在")
		}
		if userID > 0 && order.UserID != userID {
			return fmt.Errorf("订单不存在或无权操作")
		}
		if order.Status != 1 {
			return fmt.Errorf("仅待核销订单可退款")
		}
		if order.WxTransactionID == "" {
			return fmt.Errorf("订单缺少微信支付信息，无法退款")
		}
		if order.RefundNo != "" {
			return fmt.Errorf("退款处理中或已退款，请勿重复申请")
		}
		var calcErr error
		refundFen, totalFen, _, calcErr = s.CalcRefundFen(&order)
		if calcErr != nil {
			return calcErr
		}
		refundNo = fmt.Sprintf("R%s_%d", order.OrderNo, time.Now().Unix())
		res := tx.Model(&model.TicketOrder{}).
			Where("id = ? AND status = ? AND (refund_no = '' OR refund_no IS NULL)", orderID, 1).
			Updates(map[string]interface{}{
				"status":    7,
				"refund_no": refundNo,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("订单状态已变更，请刷新后重试")
		}
		return nil
	})
	if err != nil {
		return err
	}

	result, err := miniPay.CreateRefund(order.WxTransactionID, refundNo, totalFen, refundFen, "门票按付费次数剩余比例退款")
	if err != nil {
		_ = s.ReleaseRefundRequested(refundNo)
		return err
	}
	if result.RefundID != "" {
		_ = global.GVA_DB.Model(&model.TicketOrder{}).
			Where("id = ? AND refund_no = ?", order.ID, refundNo).
			Update("wx_refund_id", result.RefundID).Error
	}
	if strings.ToUpper(result.Status) == "SUCCESS" {
		return s.ApplyRefundSuccessByRefundNo(refundNo, result.RefundID, "", refundFen)
	}
	return nil
}

// ReleaseRefundRequested 微信退款失败/关闭时释放退款中状态
func (s *ticketOrder) ReleaseRefundRequested(refundNo string) error {
	if strings.TrimSpace(refundNo) == "" {
		return fmt.Errorf("缺少商户退款单号 out_refund_no")
	}
	return global.GVA_DB.Model(&model.TicketOrder{}).
		Where("refund_no = ? AND status = ?", refundNo, 7).
		Updates(map[string]interface{}{
			"status":       1,
			"refund_no":    "",
			"wx_refund_id": "",
		}).Error
}

// ApplyRefundSuccessByRefundNo 退款成功：写实退金额、按 R5 回退日历 sold
func (s *ticketOrder) ApplyRefundSuccessByRefundNo(refundNo, refundID, successTime string, refundFenHint int) error {
	if strings.TrimSpace(refundNo) == "" {
		return fmt.Errorf("缺少商户退款单号 out_refund_no")
	}
	var order model.TicketOrder
	if err := global.GVA_DB.Where("refund_no = ?", refundNo).First(&order).Error; err != nil {
		return fmt.Errorf("退款对应订单不存在")
	}

	refundAt := time.Now()
	if t, err := time.Parse(time.RFC3339, successTime); err == nil {
		refundAt = t
	}

	if order.Status == 6 {
		if order.WxRefundID != "" && refundID != "" && order.WxRefundID != refundID {
			return fmt.Errorf("微信退款单号与已退款记录不一致")
		}
		if order.WxRefundID == "" && refundID != "" {
			_ = global.GVA_DB.Model(&model.TicketOrder{}).Where("id = ?", order.ID).Update("wx_refund_id", refundID).Error
		}
		return nil
	}
	if order.Status != 1 && order.Status != 7 {
		return fmt.Errorf("订单状态不允许确认退款: status=%d", order.Status)
	}

	remainingPaid := s.RemainingPaidUseTimes(&order)
	refundFen := refundFenHint
	if refundFen <= 0 {
		rf, _, _, calcErr := s.CalcRefundFen(&order)
		if calcErr != nil {
			refundFen = 0
		} else {
			refundFen = rf
		}
	}
	refundAmount := float64(refundFen) / 100.0

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		updateRes := tx.Model(&model.TicketOrder{}).
			Where("id = ? AND status IN ? AND refund_no = ?", order.ID, []int{1, 7}, refundNo).
			Updates(map[string]interface{}{
				"status":        6,
				"wx_refund_id":  refundID,
				"refund_time":   refundAt,
				"refund_amount": refundAmount,
			})
		if updateRes.Error != nil {
			return updateRes.Error
		}
		if updateRes.RowsAffected == 0 {
			var fresh model.TicketOrder
			if err := tx.Where("id = ?", order.ID).First(&fresh).Error; err != nil {
				return fmt.Errorf("订单不存在")
			}
			if fresh.Status == 6 {
				return nil
			}
			return fmt.Errorf("订单状态已变更，请刷新后重试")
		}

		if remainingPaid <= 0 {
			return nil
		}
		// R5：多次票每单 1 张 sold-=1；单次票 sold-=remainingPaid（且不超过 quantity）
		isMulti := false
		var sku model.TicketSku
		if e := tx.Where("id = ?", order.SkuID).First(&sku).Error; e == nil {
			isMulti = sku.TicketType == 2
		} else if order.Quantity == 1 && orderPaidUseTimes(&order) > 1 {
			// SKU 缺失时：quantity=1 且 A>1 视为多次票
			isMulti = true
		}
		release := remainingPaid
		if isMulti {
			release = 1
		}
		if order.Quantity > 0 && release > order.Quantity {
			release = order.Quantity
		}
		if release <= 0 {
			return nil
		}
		calendarRes := tx.Model(&model.TicketCalendar{}).
			Where("sku_id = ? AND visit_date = ? AND sold >= ?", order.SkuID, order.VisitDate, release).
			UpdateColumn("sold", gorm.Expr("sold - ?", release))
		if calendarRes.Error != nil {
			return calendarRes.Error
		}
		if calendarRes.RowsAffected == 0 {
			global.GVA_LOG.Warn("ticket refund sold rollback skipped",
				zap.Uint("order_id", order.ID),
				zap.Uint("sku_id", order.SkuID),
				zap.Time("visit_date", order.VisitDate),
				zap.Int("release", release),
			)
		}
		return nil
	})
}

func (s *ticketOrder) GetList(req request.TicketOrderSearch) (list []model.TicketOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketOrder{})
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.BookerPhone != "" {
		db = db.Where("booker_phone LIKE ?", "%"+req.BookerPhone+"%")
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
			db = db.Where("status IN (?)", []int{1, 7}).
				Where("visit_date >= ?", today)
		case "completed", "已完成":
			db = db.Where("status IN (?)", []int{2, 3, 4, 5, 6})
		}
	}
	if req.TicketType != nil {
		db = db.Where("EXISTS (SELECT 1 FROM ticket_sku s WHERE s.id = orders.sku_id AND s.ticket_type = ? AND s.deleted_at IS NULL)", *req.TicketType)
	}
	if req.Venue != "" {
		db = db.Where("EXISTS (SELECT 1 FROM order_verify_records r WHERE r.order_id = orders.id AND r.venue = ? AND r.deleted_at IS NULL)", req.Venue)
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
	if err != nil {
		return
	}
	for i := range list {
		s.fillProductName(&list[i])
	}
	return
}

func (s *ticketOrder) GetMyList(req request.TicketOrderSearch) (list []model.TicketOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.TicketOrder{}).Where("user_deleted_at IS NULL")
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.BookerPhone != "" {
		db = db.Where("booker_phone LIKE ?", "%"+req.BookerPhone+"%")
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
			db = db.Where("status IN (?)", []int{1, 7}).
				Where("visit_date >= ?", today)
		case "completed", "已完成":
			db = db.Where("status IN (?)", []int{2, 3, 4, 5, 6})
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
	if err != nil {
		return
	}
	for i := range list {
		s.fillProductName(&list[i])
	}
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
	case 6:
		return "已退款"
	case 7:
		return "退款中"
	default:
		return "未知"
	}
}

// GetVerifyRecords 查询订单的核销记录列表
func (s *ticketOrder) GetVerifyRecords(orderID uint) (records []model.OrderVerifyRecord, err error) {
	err = global.GVA_DB.Where("order_id = ?", orderID).Order("verify_no ASC").Find(&records).Error
	return
}

// GetVenueVerifyStats 按核销月份汇总白名单场合的核销次数；month 为空取当前自然月，非法格式返回错误
func (s *ticketOrder) GetVenueVerifyStats(month string) (monthOut string, items []request.TicketVenueVerifyStatsItem, err error) {
	now := time.Now()
	var monthStart time.Time
	if strings.TrimSpace(month) == "" {
		monthStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	} else {
		monthStart, err = time.ParseInLocation("2006-01", strings.TrimSpace(month), now.Location())
		if err != nil {
			return "", nil, fmt.Errorf("月份格式错误，应为 YYYY-MM")
		}
	}
	nextMonthStart := monthStart.AddDate(0, 1, 0)
	monthOut = monthStart.Format("2006-01")

	options := model.VenueOptions()
	venueCodes := make([]string, 0, len(options))
	for _, opt := range options {
		venueCodes = append(venueCodes, opt["code"])
	}

	type venueCountRow struct {
		Venue string
		Count int64
	}
	var rows []venueCountRow
	err = global.GVA_DB.Model(&model.OrderVerifyRecord{}).
		Select("venue, COUNT(*) as count").
		Where("verified_at >= ? AND verified_at < ?", monthStart, nextMonthStart).
		Where("venue <> '' AND venue IN ?", venueCodes).
		Group("venue").
		Scan(&rows).Error
	if err != nil {
		return
	}

	countMap := make(map[string]int64, len(rows))
	for _, r := range rows {
		countMap[r.Venue] = r.Count
	}

	items = make([]request.TicketVenueVerifyStatsItem, 0, len(options))
	for _, opt := range options {
		code := opt["code"]
		items = append(items, request.TicketVenueVerifyStatsItem{
			Venue: code,
			Label: opt["label"],
			Count: countMap[code],
		})
	}
	return
}

func (s *ticketOrder) GetByID(id uint) (order model.TicketOrder, err error) {
	if err = global.GVA_DB.Where("id = ?", id).First(&order).Error; err != nil {
		return
	}
	s.fillProductName(&order)
	return
}

func (s *ticketOrder) GetMyByID(id uint, userID uint) (order model.TicketOrder, err error) {
	if err = global.GVA_DB.Where("id = ? AND user_id = ? AND user_deleted_at IS NULL", id, userID).First(&order).Error; err != nil {
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
	order.SkuTicketType = sku.TicketType
	order.SkuTicketTypeLabel = ticketTypeLabel(sku.TicketType)
}

func ticketTypeLabel(ticketType int) string {
	switch ticketType {
	case 2:
		return "多次票"
	default:
		return "单次票"
	}
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
		{
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			if visitDate.Before(today) {
				return fmt.Errorf("游玩日期不能早于今天")
			}
		}
		if req.Quantity <= 0 {
			return fmt.Errorf("购买数量必须大于 0")
		}
		if sku.TicketType == 2 && req.Quantity != 1 {
			return fmt.Errorf("多次票每单限购 1 张")
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
		m := sku.UseTimes
		if m <= 0 {
			m = 1
		}
		p := sku.GiftUseTimes
		if sku.TicketType != 2 || p < 0 {
			p = 0
		}
		paidUse := req.Quantity * m
		giftUse := req.Quantity * p
		orderNo = fmt.Sprintf("T%s%06d", time.Now().Format("20060102150405"), rand.Intn(1000000))
		order = model.TicketOrder{
			OrderNo:           orderNo,
			UserID:            userID,
			BookerName:        strings.TrimSpace(req.BookerName),
			BookerPhone:       phone,
			SkuID:             req.SkuID,
			SkuName:           sku.Name,
			Price:             sku.Price,
			Quantity:          req.Quantity,
			VisitDate:         visitDate,
			TotalAmount:       totalAmount,
			PayAmount:         totalAmount,
			Status:            0,
			PaidUseTimes:      paidUse,
			GiftUseTimes:      giftUse,
			TotalUseTimes:     paidUse + giftUse,
			VerifiedTimes:     0,
			SupportMultiVenue: sku.TicketType == 2 && sku.SupportMultiVenue,
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
// venue：支持多场合订单必填且须在白名单；否则忽略，记录为空
func (s *ticketOrder) VerifyOrder(orderID uint, venue string) error {
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

		venueCode := ""
		if order.SupportMultiVenue {
			if !model.IsValidVenue(venue) {
				return fmt.Errorf("请选择核销场合")
			}
			venueCode = venue
		}

		newVerified := order.VerifiedTimes + 1
		now := time.Now()

		record := model.OrderVerifyRecord{
			OrderID:    orderID,
			VerifyNo:   newVerified,
			VerifiedAt: now,
			Venue:      venueCode,
		}
		// 显式 Select，避免 GORM default 标签导致 venue 未写入
		if err := tx.Select("OrderID", "VerifyNo", "VerifiedAt", "Venue", "Remark").Create(&record).Error; err != nil {
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
func (s *ticketOrder) VerifyOrderByOrderNoPublic(orderNo string, venue string) error {
	var order model.TicketOrder
	if err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil || order.ID == 0 {
		return fmt.Errorf("订单不存在")
	}
	return s.VerifyOrder(order.ID, venue)
}

// DeleteMyOrder 小程序端删除本人订单，仅前台隐藏，后台仍保留
func (s *ticketOrder) DeleteMyOrder(orderID uint, userID uint) error {
	var order model.TicketOrder
	if err := global.GVA_DB.Where("id = ?", orderID).First(&order).Error; err != nil || order.ID == 0 {
		return fmt.Errorf("订单不存在")
	}
	if order.UserID != userID {
		return fmt.Errorf("无权删除该订单")
	}
	if order.UserDeletedAt != nil {
		return fmt.Errorf("订单已删除")
	}
	if order.Status != 6 {
		return fmt.Errorf("仅已退款订单可删除")
	}
	now := time.Now()
	return global.GVA_DB.Model(&model.TicketOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Update("user_deleted_at", &now).Error
}

// CloseMyPendingOrder 小程序端手动关闭本人待支付订单（仅允许未超时订单）
func (s *ticketOrder) CloseMyPendingOrder(orderID uint, userID uint, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 15 * time.Minute
	}
	deadline := time.Now().Add(-timeout)

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order model.TicketOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error; err != nil || order.ID == 0 {
			return fmt.Errorf("订单不存在")
		}
		if order.Status != 0 {
			return fmt.Errorf("仅待支付订单可关闭")
		}
		if !order.CreatedAt.After(deadline) {
			return fmt.Errorf("订单已超时，请刷新后重试")
		}

		calendarRes := tx.Model(&model.TicketCalendar{}).
			Where("sku_id = ? AND visit_date = ? AND sold >= ?", order.SkuID, order.VisitDate, order.Quantity).
			UpdateColumn("sold", gorm.Expr("sold - ?", order.Quantity))
		if calendarRes.Error != nil {
			return calendarRes.Error
		}
		if calendarRes.RowsAffected == 0 {
			return fmt.Errorf("库存回退失败，请联系客服")
		}

		res := tx.Model(&model.TicketOrder{}).
			Where("id = ? AND user_id = ? AND status = ?", orderID, userID, 0).
			Update("status", 5)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("订单状态已变更，请刷新后重试")
		}
		return nil
	})
}

// CloseTimeoutUnpaidOrders 关闭超时未支付订单（status=0 -> 5），并回退对应日历已售库存
func (s *ticketOrder) CloseTimeoutUnpaidOrders(timeout time.Duration, batchSize int) (int, error) {
	if timeout <= 0 {
		timeout = 15 * time.Minute
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	deadline := time.Now().Add(-timeout)

	var pending []model.TicketOrder
	if err := global.GVA_DB.
		Where("status = ? AND created_at <= ?", 0, deadline).
		Order("id ASC").
		Limit(batchSize).
		Find(&pending).Error; err != nil {
		return 0, err
	}
	if len(pending) == 0 {
		return 0, nil
	}

	closedCount := 0
	for _, item := range pending {
		err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			var order model.TicketOrder
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", item.ID).
				First(&order).Error; err != nil {
				return err
			}
			// 仅处理仍是待支付且已超时的订单，避免和支付回调并发冲突
			if order.Status != 0 || order.CreatedAt.After(deadline) {
				return nil
			}

			calendarRes := tx.Model(&model.TicketCalendar{}).
				Where("sku_id = ? AND visit_date = ? AND sold >= ?", order.SkuID, order.VisitDate, order.Quantity).
				UpdateColumn("sold", gorm.Expr("sold - ?", order.Quantity))
			if calendarRes.Error != nil {
				return calendarRes.Error
			}
			if calendarRes.RowsAffected == 0 {
				// 部分历史数据可能已被人工调整或库存行缺失；此时不阻断关单，避免定时任务反复卡在同一订单。
				global.GVA_LOG.Warn("ticket timeout close inventory rollback skipped",
					zap.Uint("order_id", order.ID),
					zap.Uint("sku_id", order.SkuID),
					zap.Time("visit_date", order.VisitDate),
					zap.Int("quantity", order.Quantity),
				)
			}

			res := tx.Model(&model.TicketOrder{}).
				Where("id = ? AND status = ?", order.ID, 0).
				Update("status", 5)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return nil
			}

			closedCount++
			return nil
		})
		if err != nil {
			return closedCount, err
		}
	}
	return closedCount, nil
}
