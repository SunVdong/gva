package service

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
	miniPay "github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type activityOrder struct{}

// validChineseMobile 中国大陆 11 位手机号
var validChineseMobile = regexp.MustCompile(`^1[3-9]\d{9}$`)

func (s *activityOrder) GetList(req request.ActivityOrderSearch) (list []model.ActivityOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.ActivityOrder{})
	db = s.applySearch(db, req)
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

func (s *activityOrder) GetMyList(req request.ActivityOrderSearch) (list []model.ActivityOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.ActivityOrder{})
	db = s.applySearch(db, req)
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
	s.fillActivityExtras(list)
	return
}

func (s *activityOrder) applySearch(db *gorm.DB, req request.ActivityOrderSearch) *gorm.DB {
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.ContactPhone != "" {
		db = db.Where("contact_phone LIKE ?", "%"+req.ContactPhone+"%")
	}
	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.ActivityID > 0 {
		db = db.Where("activity_id = ?", req.ActivityID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.OrderType != nil {
		switch *req.OrderType {
		case "pending_payment", "待支付":
			db = db.Where("status = ?", 0)
		case "pending_verify", "待核销":
			db = db.Where("status IN ?", []int{1, 7})
		case "completed", "已完成":
			db = db.Where("status IN ?", []int{2, 5, 6})
		}
	}
	return db
}

func (s *activityOrder) GetByID(id uint) (order model.ActivityOrder, err error) {
	if err = global.GVA_DB.Where("id = ?", id).First(&order).Error; err != nil {
		return
	}
	s.fillActivityExtra(&order)
	return
}

func (s *activityOrder) GetMyByID(id uint, userID uint) (order model.ActivityOrder, err error) {
	if err = global.GVA_DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		return
	}
	s.fillActivityExtra(&order)
	return
}

func (s *activityOrder) GetVerifyRecords(orderID uint) (records []model.ActivityOrderVerifyRecord, err error) {
	err = global.GVA_DB.Where("order_id = ?", orderID).Order("verify_no ASC").Find(&records).Error
	return
}

func (s *activityOrder) OrderStatusLabel(order *model.ActivityOrder) string {
	switch order.Status {
	case 0:
		return "待支付"
	case 1:
		if order.VerifiedTimes > 0 {
			return "核销中"
		}
		return "待核销"
	case 2:
		return "已核销"
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

func (s *activityOrder) fillActivityExtra(order *model.ActivityOrder) {
	if order.ActivityID == 0 {
		return
	}
	var act model.Activity
	if e := global.GVA_DB.Select("group_qr, service_qr, cover_image").Where("id = ?", order.ActivityID).First(&act).Error; e == nil {
		order.GroupQr = act.GroupQr
		order.ServiceQr = act.ServiceQr
		order.CoverImage = act.CoverImage
	}
}

func (s *activityOrder) fillActivityExtras(list []model.ActivityOrder) {
	ids := make([]uint, 0, len(list))
	seen := map[uint]struct{}{}
	for _, o := range list {
		if o.ActivityID == 0 {
			continue
		}
		if _, ok := seen[o.ActivityID]; ok {
			continue
		}
		seen[o.ActivityID] = struct{}{}
		ids = append(ids, o.ActivityID)
	}
	if len(ids) == 0 {
		return
	}
	var acts []model.Activity
	if err := global.GVA_DB.Select("id, group_qr, service_qr, cover_image").Where("id IN ?", ids).Find(&acts).Error; err != nil {
		return
	}
	m := make(map[uint]model.Activity, len(acts))
	for _, a := range acts {
		m[a.ID] = a
	}
	for i := range list {
		if a, ok := m[list[i].ActivityID]; ok {
			list[i].GroupQr = a.GroupQr
			list[i].ServiceQr = a.ServiceQr
			list[i].CoverImage = a.CoverImage
		}
	}
}

// CreateOrder 小程序报名下单：校验活动时间窗/显示/名额，事务占用 sold
func (s *activityOrder) CreateOrder(userID uint, req request.MiniOrderCreate) (order model.ActivityOrder, err error) {
	phone := strings.TrimSpace(req.ContactPhone)
	if !validChineseMobile.MatchString(phone) {
		return order, fmt.Errorf("联系人手机号格式不正确")
	}
	name := strings.TrimSpace(req.ContactName)
	if name == "" {
		return order, fmt.Errorf("联系人姓名不能为空")
	}
	if req.Quantity <= 0 {
		return order, fmt.Errorf("人次必须大于0")
	}

	var orderNo string
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var act model.Activity
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", req.ActivityID).First(&act).Error; e != nil {
			return fmt.Errorf("活动不存在")
		}
		if act.Status != 1 {
			return fmt.Errorf("活动未开放报名")
		}
		now := time.Now()
		if now.Before(act.StartTime) {
			return fmt.Errorf("活动尚未开始，无法报名")
		}
		if now.After(act.EndTime) {
			return fmt.Errorf("活动已结束，无法报名")
		}
		ur := tx.Model(&model.Activity{}).
			Where("id = ? AND sold + ? <= quota", act.ID, req.Quantity).
			UpdateColumn("sold", gorm.Expr("sold + ?", req.Quantity))
		if ur.Error != nil {
			return ur.Error
		}
		if ur.RowsAffected == 0 {
			return fmt.Errorf("名额不足")
		}

		payAmount := act.Price * float64(req.Quantity)
		if payAmount <= 0 || act.Price <= 0 {
			return fmt.Errorf("活动价格异常，无法报名")
		}
		orderNo = fmt.Sprintf("A%s%06d", time.Now().Format("20060102150405"), rand.Intn(1000000))
		order = model.ActivityOrder{
			OrderNo:       orderNo,
			UserID:        userID,
			ActivityID:    act.ID,
			ActivityName:  act.Name,
			ContactName:   name,
			ContactPhone:  phone,
			Quantity:      req.Quantity,
			UnitPrice:     act.Price,
			PayAmount:     payAmount,
			Status:        0,
			TotalUseTimes: req.Quantity,
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

// GetByOrderNoPublic 公开：按订单号查询
func (s *activityOrder) GetByOrderNoPublic(orderNo string) (order model.ActivityOrder, err error) {
	if err = global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return
	}
	s.fillActivityExtra(&order)
	return
}

// VerifyOrder 核销一次（人次 +1）
func (s *activityOrder) VerifyOrder(orderID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order model.ActivityOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).First(&order).Error; err != nil || order.ID == 0 {
			return fmt.Errorf("订单不存在")
		}
		if order.Status == 2 {
			return fmt.Errorf("该订单已核销完毕")
		}
		if order.Status == 6 || order.Status == 7 {
			return fmt.Errorf("订单已退款或退款中，不可核销")
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
		record := model.ActivityOrderVerifyRecord{
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
		return tx.Model(&model.ActivityOrder{}).Where("id = ?", orderID).Updates(updates).Error
	})
}

func (s *activityOrder) VerifyOrderByOrderNoPublic(orderNo string) error {
	var order model.ActivityOrder
	if err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil || order.ID == 0 {
		return fmt.Errorf("订单不存在")
	}
	return s.VerifyOrder(order.ID)
}

// CalcRefundFen 按未核销比例计算退款分
func (s *activityOrder) CalcRefundFen(order *model.ActivityOrder) (refundFen, totalFen, remaining int, err error) {
	totalUse := order.TotalUseTimes
	if totalUse <= 0 {
		return 0, 0, 0, fmt.Errorf("订单可核销次数异常")
	}
	remaining = totalUse - order.VerifiedTimes
	if remaining <= 0 {
		return 0, 0, 0, fmt.Errorf("已全部核销，不可退款")
	}
	totalFen = int(math.Round(order.PayAmount * 100))
	if totalFen <= 0 {
		return 0, 0, 0, fmt.Errorf("订单金额异常")
	}
	refundFen = int(math.Round(order.PayAmount * 100 * float64(remaining) / float64(totalUse)))
	if refundFen <= 0 {
		return 0, 0, 0, fmt.Errorf("退款金额异常")
	}
	if refundFen > totalFen {
		refundFen = totalFen
	}
	return
}

// AdminRefund 后台按未核销比例发起微信退款。
// 先落库 status=7 再调微信，避免退款受理窗口内仍可核销导致比例不一致。
func (s *activityOrder) AdminRefund(orderID uint) error {
	var order model.ActivityOrder
	var refundFen, totalFen int
	var refundNo string

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderID).First(&order).Error; e != nil || order.ID == 0 {
			return fmt.Errorf("订单不存在")
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
		res := tx.Model(&model.ActivityOrder{}).
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

	result, err := miniPay.CreateRefund(order.WxTransactionID, refundNo, totalFen, refundFen, "限时活动按未核销比例退款")
	if err != nil {
		_ = s.ReleaseRefundRequested(refundNo)
		return err
	}
	if result.RefundID != "" {
		_ = global.GVA_DB.Model(&model.ActivityOrder{}).
			Where("id = ? AND refund_no = ?", order.ID, refundNo).
			Update("wx_refund_id", result.RefundID).Error
	}
	if strings.ToUpper(result.Status) == "SUCCESS" {
		return s.ApplyRefundSuccessByRefundNo(refundNo, result.RefundID, "", refundFen)
	}
	return nil
}

func (s *activityOrder) ReleaseRefundRequested(refundNo string) error {
	if strings.TrimSpace(refundNo) == "" {
		return fmt.Errorf("缺少商户退款单号 out_refund_no")
	}
	return global.GVA_DB.Model(&model.ActivityOrder{}).
		Where("refund_no = ? AND status = ?", refundNo, 7).
		Updates(map[string]interface{}{
			"status":       1,
			"refund_no":    "",
			"wx_refund_id": "",
		}).Error
}

// ApplyRefundSuccessByRefundNo 退款成功：写实退金额、回退未核销名额、禁止再核销
func (s *activityOrder) ApplyRefundSuccessByRefundNo(refundNo, refundID, successTime string, refundFenHint int) error {
	if strings.TrimSpace(refundNo) == "" {
		return fmt.Errorf("缺少商户退款单号 out_refund_no")
	}
	var order model.ActivityOrder
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
			_ = global.GVA_DB.Model(&model.ActivityOrder{}).Where("id = ?", order.ID).Update("wx_refund_id", refundID).Error
		}
		return nil
	}
	if order.Status != 1 && order.Status != 7 {
		return fmt.Errorf("订单状态不允许确认退款: status=%d", order.Status)
	}

	remaining := order.TotalUseTimes - order.VerifiedTimes
	if remaining < 0 {
		remaining = 0
	}
	refundFen := refundFenHint
	if refundFen <= 0 {
		rf, _, _, calcErr := s.CalcRefundFen(&order)
		if calcErr != nil {
			// 极端：已全核销但仍在退款中，按 0 处理金额但应已不可能走到这里
			refundFen = 0
		} else {
			refundFen = rf
		}
	}
	refundAmount := float64(refundFen) / 100.0

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		updateRes := tx.Model(&model.ActivityOrder{}).
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
			var fresh model.ActivityOrder
			if err := tx.Where("id = ?", order.ID).First(&fresh).Error; err != nil {
				return fmt.Errorf("订单不存在")
			}
			if fresh.Status == 6 {
				return nil
			}
			return fmt.Errorf("订单状态已变更，请刷新后重试")
		}
		if remaining > 0 {
			actRes := tx.Model(&model.Activity{}).
				Where("id = ? AND sold >= ?", order.ActivityID, remaining).
				UpdateColumn("sold", gorm.Expr("sold - ?", remaining))
			if actRes.Error != nil {
				return actRes.Error
			}
			if actRes.RowsAffected == 0 {
				global.GVA_LOG.Warn("limitedActivity refund sold rollback skipped",
					zap.Uint("order_id", order.ID),
					zap.Uint("activity_id", order.ActivityID),
					zap.Int("remaining", remaining),
				)
			}
		}
		return nil
	})
}

// CloseTimeoutUnpaidOrders 关闭超时未支付订单并释放名额
func (s *activityOrder) CloseTimeoutUnpaidOrders(timeout time.Duration, batchSize int) (int, error) {
	if timeout <= 0 {
		timeout = 15 * time.Minute
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	deadline := time.Now().Add(-timeout)

	var pending []model.ActivityOrder
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
			var order model.ActivityOrder
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", item.ID).First(&order).Error; err != nil {
				return err
			}
			if order.Status != 0 || order.CreatedAt.After(deadline) {
				return nil
			}
			actRes := tx.Model(&model.Activity{}).
				Where("id = ? AND sold >= ?", order.ActivityID, order.Quantity).
				UpdateColumn("sold", gorm.Expr("sold - ?", order.Quantity))
			if actRes.Error != nil {
				return actRes.Error
			}
			if actRes.RowsAffected == 0 {
				global.GVA_LOG.Warn("limitedActivity timeout close sold rollback skipped",
					zap.Uint("order_id", order.ID),
					zap.Uint("activity_id", order.ActivityID),
					zap.Int("quantity", order.Quantity),
				)
			}
			res := tx.Model(&model.ActivityOrder{}).
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
