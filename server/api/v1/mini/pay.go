package mini

import (
	"io"
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
// @Description 根据订单类型与订单 ID 生成预支付单，返回小程序调起支付所需参数（需登录，Header 带 x-token）
// @Accept      json
// @Produce     json
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
		response.OkWithData(params, c)
		return
	default:
		response.FailWithMessage("不支持的订单类型", c)
	}
}

// Notify 微信支付结果回调（由微信服务器调用，不展示在接口文档中）
func (a *PayApi) Notify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.XML(200, gin.H{"return_code": "FAIL", "return_msg": "读取body失败"})
		return
	}
	result, err := mini.ParseAndVerifyPaidNotify(body)
	if err != nil {
		c.XML(200, gin.H{"return_code": "FAIL", "return_msg": err.Error()})
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
	c.Header("Content-Type", "application/xml")
	c.Status(200)
	if err := mini.RespondPaidNotifySuccess(c.Writer); err != nil {
		c.Writer.WriteString("<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[响应失败]]></return_msg></xml>")
	}
}
