package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var Order = new(miniOrderApi)

type miniOrderApi struct{}

// Create 小程序-报名下单
// @Tags        小程序-限时活动
// @Summary     报名创建订单
// @Description 提交人次、联系人姓名与手机号创建待支付订单并占用名额。支付请调用 POST /mini/pay/create，body 传 {"orderType":"limitedActivity","orderId":订单ID}
// @Accept      json
// @Produce     json
// @Param       x-token header string true "小程序登录 token"
// @Param       data body request.MiniOrderCreate true "报名信息"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /limitedActivity/mini/order/create [post]
func (a *miniOrderApi) Create(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req request.MiniOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, err := svcOrder.CreateOrder(userID, req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(order, c)
}

// MyList 小程序-我的活动订单
// @Tags        小程序-限时活动
// @Summary     我的活动订单列表
// @Description orderType: pending_payment|pending_verify|completed；不传返回全部
// @Accept      json
// @Produce     json
// @Param       x-token header string true "小程序登录 token"
// @Param       orderType query string false "pending_payment|pending_verify|completed"
// @Param       page query int false "页码"
// @Param       pageSize query int false "每页条数"
// @Success     200 {object} response.Response{data=response.PageResult,msg=string}
// @Router      /limitedActivity/mini/order/myList [get]
func (a *miniOrderApi) MyList(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req request.ActivityOrderSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = userID
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}
	list, total, err := svcOrder.GetMyList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	items := make([]gin.H, 0, len(list))
	for _, o := range list {
		totalUse := o.TotalUseTimes
		if totalUse <= 0 {
			totalUse = 1
		}
		remaining := totalUse - o.VerifiedTimes
		if o.Status == 6 || o.Status == 7 {
			remaining = 0
		}
		if remaining < 0 {
			remaining = 0
		}
		items = append(items, gin.H{
			"id":             o.ID,
			"orderNo":        o.OrderNo,
			"userId":         o.UserID,
			"activityId":     o.ActivityID,
			"activityName":   o.ActivityName,
			"contactName":    o.ContactName,
			"contactPhone":   o.ContactPhone,
			"quantity":       o.Quantity,
			"unitPrice":      o.UnitPrice,
			"payAmount":      o.PayAmount,
			"status":         o.Status,
			"totalUseTimes":  totalUse,
			"verifiedTimes":  o.VerifiedTimes,
			"remainingTimes": remaining,
			"payTime":        o.PayTime,
			"createdAt":      o.CreatedAt,
			"statusLabel":    svcOrder.OrderStatusLabel(&o),
			"coverImage":     o.CoverImage,
			"longImage":      o.LongImage,
			"groupQr":        o.GroupQr,
			"serviceQr":      o.ServiceQr,
		})
	}
	response.OkWithDetailed(response.PageResult{
		List: items, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Detail 小程序-订单详情
// @Tags        小程序-限时活动
// @Summary     活动订单详情
// @Description 含核销进度、群二维码、客服二维码；用户不可自助退款
// @Accept      json
// @Produce     json
// @Param       x-token header string true "小程序登录 token"
// @Param       id query int true "订单ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /limitedActivity/mini/order/detail [get]
func (a *miniOrderApi) Detail(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order, err := svcOrder.GetMyByID(idReq.ID, userID)
	if err != nil || order.ID == 0 {
		response.FailWithMessage("订单不存在", c)
		return
	}
	totalUse := order.TotalUseTimes
	if totalUse <= 0 {
		totalUse = 1
	}
	remaining := totalUse - order.VerifiedTimes
	if order.Status == 6 || order.Status == 7 {
		remaining = 0
	}
	if remaining < 0 {
		remaining = 0
	}
	verifyRecords, _ := svcOrder.GetVerifyRecords(order.ID)
	response.OkWithData(gin.H{
		"order":          order,
		"remainingTimes": remaining,
		"verifyRecords":  verifyRecords,
		"canRefund":      false, // 活动订单不支持用户自助退款，请联系客服
		"statusLabel":    svcOrder.OrderStatusLabel(&order),
	}, c)
}
