package mini

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	userModel "github.com/flipped-aurora/gin-vue-admin/server/model/user"
	ticketModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PayApi struct{}

// Create 调起微信支付（JSAPI），获取小程序 wx.requestPayment 所需参数
// @Tags        小程序
// @Summary     调起支付
// @Description 根据订单类型与订单 ID 生成预支付单（微信 V3），返回小程序调起支付所需参数（signType 为 RSA）。需登录，请求头必带 x-token；须已完整配置微信支付，否则返回错误。
// @Accept      json
// @Produce     json
// @Param       x-token header string true "小程序登录后返回的 token（必填）"
// @Param       data body object true "请求体" example({"orderType":"ticket","orderId":1})
// @Success     200 {object} response.Response{data=object,msg=string} "data 含 timeStamp,nonceStr,package,signType,paySign"
// @Router      /mini/pay/create [post]
func (a *PayApi) Create(c *gin.Context) {
	userIDVal, exists := c.Get("x-user-id")
	if !exists || userIDVal == nil {
		response.FailWithMessage("请先登录", c)
		return
	}
	userID, _ := userIDVal.(uint)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}

	var req struct {
		OrderType string `json:"orderType" binding:"required"` // 订单类型：ticket 景点门票
		OrderID   uint   `json:"orderId" binding:"required"`   // 订单 ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var openID string
	{
		var u userModel.User
		if err := global.GVA_DB.Select("openid").Where("id = ?", userID).First(&u).Error; err != nil || u.OpenID == "" {
			response.FailWithMessage("用户未绑定微信，无法支付", c)
			return
		}
		openID = u.OpenID
	}

	switch req.OrderType {
	case "ticket":
		var order ticketModel.TicketOrder
		if err := global.GVA_DB.Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
			response.FailWithMessage("订单不存在或无权支付", c)
			return
		}
		if order.Status != 0 {
			response.FailWithMessage("订单状态不允许支付", c)
			return
		}
		// 金额转为分，微信单位是分；须与 ticketPayNotifyAssertAmountAndTx 中 math.Round 一致，避免 float 误差导致下单与回调金额不一致
		fen := int64(math.Round(order.PayAmount * 100))
		if fen <= 0 {
			response.FailWithMessage("订单金额异常", c)
			return
		}
		// 每次发起支付使用不同的 out_trade_no，避免微信预支付单过期后同号不可复用
		outTradeNo := fmt.Sprintf("%s_%d", order.OrderNo, time.Now().Unix())
		params, err := mini.CreateJSAPI(outTradeNo, fen, "景点门票-"+order.OrderNo, openID, c.ClientIP())
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithData(params, c)
		return
	default:
		response.FailWithMessage("不支持的订单类型", c)
	}
}

// Notify 微信支付 V3 结果回调（由微信服务器 POST JSON，不展示在接口文档中）
func (a *PayApi) Notify(c *gin.Context) {
	result, err := mini.ParseAndVerifyPaidNotify(c.Request)
	if err != nil {
		c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}
	// 从 out_trade_no 提取真实订单号（pay/create 追加了 _unix 后缀以支持重复发起支付）
	orderNo := result.OutTradeNo
	if idx := strings.Index(orderNo, "_"); idx > 0 {
		orderNo = orderNo[:idx]
	}
	// 根据订单号前缀区分业务：T=门票
	if len(orderNo) >= 1 && orderNo[0] == 'T' {
		if err := applyTicketOrderPayNotify(orderNo, result); err != nil {
			c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
			return
		}
	} else {
		c.JSON(200, gin.H{"code": "FAIL", "message": "未知的订单类型: " + result.OutTradeNo})
		return
	}
	c.JSON(200, gin.H{"code": "SUCCESS", "message": "成功"})
}

// applyTicketOrderPayNotify 验金额、微信订单号，更新或幂等；依赖 wx_transaction_id 区分「同一笔支付重复通知」与「不同支付」。
func applyTicketOrderPayNotify(orderNo string, result *mini.PaidNotifyResult) error {
	if result.TotalFee <= 0 {
		return fmt.Errorf("回调金额无效")
	}
	if result.TransactionID == "" {
		return fmt.Errorf("缺少微信订单号 transaction_id")
	}
	var order ticketModel.TicketOrder
	if err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}
	if err := ticketPayNotifyAssertAmountAndTx(&order, result); err != nil {
		return err
	}
	switch order.Status {
	case 1:
		return ticketPayNotifyIdempotentPaid(&order, result)
	case 0:
		now := time.Now()
		res := global.GVA_DB.Model(&ticketModel.TicketOrder{}).
			Where("order_no = ? AND status = ?", orderNo, 0).
			Updates(map[string]interface{}{
				"status":            1,
				"pay_time":          now,
				"wx_transaction_id": result.TransactionID,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			return nil
		}
		// 并发：另一请求已把订单置为已支付
		if err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			return fmt.Errorf("订单不存在")
		}
		return ticketPayNotifyIdempotentPaid(&order, result)
	default:
		return fmt.Errorf("订单状态不允许确认支付: status=%d", order.Status)
	}
}

func ticketPayNotifyAssertAmountAndTx(order *ticketModel.TicketOrder, result *mini.PaidNotifyResult) error {
	expectedFen := int(math.Round(order.PayAmount * 100))
	if result.TotalFee != expectedFen {
		return fmt.Errorf("支付金额与订单不一致: 订单应付%d分, 通知%d分", expectedFen, result.TotalFee)
	}
	if order.WxTransactionID != "" && order.WxTransactionID != result.TransactionID {
		return fmt.Errorf("微信订单号与已支付记录不一致")
	}
	return nil
}

// ticketPayNotifyIdempotentPaid 订单已为已支付：仅允许同一 transaction_id（或补写历史空字段）的重复通知。
func ticketPayNotifyIdempotentPaid(order *ticketModel.TicketOrder, result *mini.PaidNotifyResult) error {
	if order.Status != 1 {
		return fmt.Errorf("订单状态异常: status=%d", order.Status)
	}
	if err := ticketPayNotifyAssertAmountAndTx(order, result); err != nil {
		return err
	}
	if order.WxTransactionID == result.TransactionID {
		return nil
	}
	if order.WxTransactionID != "" {
		return fmt.Errorf("微信订单号与已支付记录不一致")
	}
	res := global.GVA_DB.Model(&ticketModel.TicketOrder{}).
		Where("order_no = ? AND status = ? AND wx_transaction_id = ?", order.OrderNo, 1, "").
		Update("wx_transaction_id", result.TransactionID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	var fresh ticketModel.TicketOrder
	if err := global.GVA_DB.Where("order_no = ?", order.OrderNo).First(&fresh).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}
	if fresh.WxTransactionID != result.TransactionID {
		return fmt.Errorf("微信订单号与已支付记录不一致")
	}
	return nil
}

// Refund 申请退款（仅待核销且 verified_times=0 可退；多次票部分核销后不可退，全额退款）
// @Tags        小程序
// @Summary     申请退款
// @Description 对已支付且待核销、且尚未产生核销次数的门票订单申请全额退款（多次票若已核销过任一次则不可退）。需登录，请求头必带 x-token。
// @Accept      json
// @Produce     json
// @Param       x-token header string true "小程序登录后返回的 token（必填）"
// @Param       data body object true "请求体" example({"orderId":1})
// @Success     200 {object} response.Response{msg=string}
// @Router      /mini/pay/refund [post]
func (a *PayApi) Refund(c *gin.Context) {
	userIDVal, exists := c.Get("x-user-id")
	if !exists || userIDVal == nil {
		response.FailWithMessage("请先登录", c)
		return
	}
	userID, _ := userIDVal.(uint)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}

	var req struct {
		OrderID uint `json:"orderId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var order ticketModel.TicketOrder
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		response.FailWithMessage("订单不存在或无权操作", c)
		return
	}
	if order.Status != 1 {
		response.FailWithMessage("当前订单状态不允许退款", c)
		return
	}
	if order.VerifiedTimes > 0 {
		response.FailWithMessage("订单已产生核销记录（含多次票已使用次数），不可退款", c)
		return
	}
	if order.WxTransactionID == "" {
		response.FailWithMessage("订单缺少微信支付信息，无法退款", c)
		return
	}
	if order.RefundNo != "" {
		response.FailWithMessage("退款处理中或已退款，请勿重复申请", c)
		return
	}

	totalFen := int(math.Round(order.PayAmount * 100))
	if totalFen <= 0 {
		response.FailWithMessage("订单金额异常", c)
		return
	}
	refundNo := fmt.Sprintf("R%s_%d", order.OrderNo, time.Now().Unix())
	result, err := mini.CreateRefund(order.WxTransactionID, refundNo, totalFen, totalFen, "用户申请退款")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ticketRefundMarkRequested(order.ID, refundNo, result.RefundID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.ToUpper(result.Status) != "SUCCESS" {
		response.OkWithMessage("退款申请已受理，请稍后查看结果", c)
		return
	}
	if err := applyTicketOrderRefundSuccessByRefundNo(refundNo, result.RefundID, ""); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("退款成功", c)
}

// RefundNotify 微信退款结果回调（由微信服务器 POST JSON，不展示在接口文档中）
func (a *PayApi) RefundNotify(c *gin.Context) {
	result, err := mini.ParseAndVerifyRefundNotify(c.Request)
	if err != nil {
		c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}
	status := strings.ToUpper(strings.TrimSpace(result.RefundStatus))
	switch status {
	case "SUCCESS":
		if err := applyTicketOrderRefundSuccessByRefundNo(result.OutRefundNo, result.RefundID, result.SuccessTime); err != nil {
			c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
			return
		}
	case "CLOSED", "ABNORMAL":
		if err := ticketRefundReleaseRequested(result.OutRefundNo); err != nil {
			c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
			return
		}
	default:
		// PROCESSING 等中间态，保持已受理状态，等待后续通知
	}
	c.JSON(200, gin.H{"code": "SUCCESS", "message": "成功"})
}

func ticketRefundMarkRequested(orderID uint, refundNo, wxRefundID string) error {
	res := global.GVA_DB.Model(&ticketModel.TicketOrder{}).
		Where("id = ? AND status = ? AND (refund_no = '' OR refund_no IS NULL OR refund_no = ?)", orderID, 1, refundNo).
		Updates(map[string]interface{}{
			"status":       7,
			"refund_no":    refundNo,
			"wx_refund_id": wxRefundID,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}

	var fresh ticketModel.TicketOrder
	if err := global.GVA_DB.Where("id = ?", orderID).First(&fresh).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}
	if fresh.Status == 6 {
		return fmt.Errorf("订单已退款")
	}
	if fresh.Status == 7 && fresh.RefundNo == refundNo {
		return nil
	}
	if fresh.RefundNo == refundNo {
		return nil
	}
	if fresh.RefundNo != "" {
		return fmt.Errorf("退款处理中或已退款，请勿重复申请")
	}
	return fmt.Errorf("订单状态已变更，请刷新后重试")
}

func ticketRefundReleaseRequested(refundNo string) error {
	if strings.TrimSpace(refundNo) == "" {
		return fmt.Errorf("缺少商户退款单号 out_refund_no")
	}
	res := global.GVA_DB.Model(&ticketModel.TicketOrder{}).
		Where("refund_no = ? AND status = ?", refundNo, 7).
		Updates(map[string]interface{}{
			"status":       1,
			"refund_no":    "",
			"wx_refund_id": "",
		})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func applyTicketOrderRefundSuccessByRefundNo(refundNo, refundID, successTime string) error {
	if strings.TrimSpace(refundNo) == "" {
		return fmt.Errorf("缺少商户退款单号 out_refund_no")
	}
	var order ticketModel.TicketOrder
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
			_ = global.GVA_DB.Model(&ticketModel.TicketOrder{}).Where("id = ?", order.ID).Update("wx_refund_id", refundID).Error
		}
		return nil
	}
	if order.Status != 1 && order.Status != 7 {
		return fmt.Errorf("订单状态不允许确认退款: status=%d", order.Status)
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		updateRes := tx.Model(&ticketModel.TicketOrder{}).
			Where("id = ? AND status IN (?) AND refund_no = ?", order.ID, []int{1, 7}, refundNo).
			Updates(map[string]interface{}{
				"status":       6,
				"wx_refund_id": refundID,
				"refund_time":  refundAt,
			})
		if updateRes.Error != nil {
			return updateRes.Error
		}
		if updateRes.RowsAffected == 0 {
			var fresh ticketModel.TicketOrder
			if err := tx.Where("id = ?", order.ID).First(&fresh).Error; err != nil {
				return fmt.Errorf("订单不存在")
			}
			if fresh.Status == 6 {
				if fresh.WxRefundID != "" && refundID != "" && fresh.WxRefundID != refundID {
					return fmt.Errorf("微信退款单号与已退款记录不一致")
				}
				return nil
			}
			if fresh.Status == 7 {
				return fmt.Errorf("退款结果处理中，请稍后重试")
			}
			return fmt.Errorf("订单状态已变更，请刷新后重试")
		}

		calendarRes := tx.Model(&ticketModel.TicketCalendar{}).
			Where("sku_id = ? AND visit_date = ? AND sold >= ?", order.SkuID, order.VisitDate, order.Quantity).
			UpdateColumn("sold", gorm.Expr("sold - ?", order.Quantity))
		if calendarRes.Error != nil {
			return calendarRes.Error
		}
		if calendarRes.RowsAffected == 0 {
			return fmt.Errorf("退款成功但回退库存失败，请联系客服")
		}
		return nil
	})
}
