package mini

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	userModel "github.com/flipped-aurora/gin-vue-admin/server/model/user"
	ticketModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"github.com/gin-gonic/gin"
)

type PayApi struct{}

// Create 调起微信支付（JSAPI），获取小程序 wx.requestPayment 所需参数
// @Tags        小程序
// @Summary     调起支付
// @Description 根据订单类型与订单 ID 生成预支付单（微信 V3），返回小程序调起支付所需参数（signType 为 RSA）。需登录，请求头必带 x-token。未配置完整（app-id、mch-id、api-v3-key、mch-api-serial-no、notify-url、商户私钥）时返回模拟参数（data.mock=true），订单会直接置为已支付，便于联调；前端可根据 data.mock 或 data.paySign==\"MOCK_SIMULATION\" 跳过 wx.requestPayment 并提示模拟成功。
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
		// 金额转为分，微信单位是分
		fen := int64(order.PayAmount * 100)
		if fen <= 0 {
			response.FailWithMessage("订单金额异常", c)
			return
		}
		params, err := mini.CreateJSAPI(order.OrderNo, fen, "景点门票-"+order.OrderNo, openID, c.ClientIP())
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		// 模拟支付：未配置微信支付时返回 mock 参数，并直接将订单置为已支付便于联调
		if params.PaySign == mini.MockPaySign {
			global.GVA_DB.Model(&ticketModel.TicketOrder{}).Where("id = ?", order.ID).
				Updates(map[string]interface{}{"status": 1, "pay_time": time.Now()})
			response.OkWithData(gin.H{"timeStamp": params.TimeStamp, "nonceStr": params.NonceStr, "package": params.Package, "signType": params.SignType, "paySign": params.PaySign, "mock": true}, c)
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
	// 根据商户订单号更新业务订单（订单号前缀区分业务：T=门票）
	outTradeNo := result.OutTradeNo
	if len(outTradeNo) >= 1 && outTradeNo[0] == 'T' {
		now := time.Now()
		global.GVA_DB.Model(&ticketModel.TicketOrder{}).
			Where("order_no = ?", outTradeNo).
			Updates(map[string]interface{}{
				"status":   1,
				"pay_time": now,
			})
	}
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(200)
	if err := mini.RespondPaidNotifySuccess(c.Writer); err != nil {
		_, _ = c.Writer.Write([]byte(`{"code":"FAIL","message":"响应失败"}`))
	}
}
