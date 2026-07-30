package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	laRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
	"github.com/gin-gonic/gin"
)

var Order = new(activityOrderApi)

type activityOrderApi struct{}

// GetList 后台订单列表
// @Tags LimitedActivityOrder
// @Summary 活动订单列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param data query laRequest.ActivityOrderSearch true "分页搜索"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /limitedActivity/order/getOrderList [get]
func (a *activityOrderApi) GetList(c *gin.Context) {
	var req laRequest.ActivityOrderSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceOrder.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Find 后台订单详情
// @Tags LimitedActivityOrder
// @Summary 活动订单详情
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query int true "订单ID"
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /limitedActivity/order/findOrder [get]
func (a *activityOrderApi) Find(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, err := serviceOrder.GetByID(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	totalUse := order.TotalUseTimes
	if totalUse <= 0 {
		totalUse = 1
	}
	remaining := totalUse - order.VerifiedTimes
	if remaining < 0 {
		remaining = 0
	}
	verifyRecords, _ := serviceOrder.GetVerifyRecords(order.ID)
	canRefund := order.Status == 1 && remaining > 0 && order.WxTransactionID != ""
	refundFen := 0
	if canRefund {
		if rf, _, _, calcErr := serviceOrder.CalcRefundFen(&order); calcErr == nil {
			refundFen = rf
		} else {
			canRefund = false
		}
	}
	data := gin.H{
		"order":           order,
		"remainingTimes":  remaining,
		"verifyRecords":   verifyRecords,
		"canRefund":       canRefund,
		"refundAmountFen": refundFen,
		"refundAmount":    float64(refundFen) / 100.0,
	}
	if order.Status == 2 && order.VerifiedAt != nil {
		review, _ := serviceOrderReview.GetByOrderID(order.ID)
		if review.ID != 0 {
			data["review"] = gin.H{
				"ID":        review.ID,
				"rating":    review.Rating,
				"content":   review.Content,
				"createdAt": review.CreatedAt,
			}
		} else {
			data["review"] = nil
		}
	}
	response.OkWithData(data, c)
}

// Refund 后台按未核销比例退款
// @Tags LimitedActivityOrder
// @Summary 活动订单退款(按未核销比例)
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query int true "订单ID"
// @Success 200 {object} response.Response{msg=string} "退款成功或已受理"
// @Router /limitedActivity/order/refundOrder [post]
func (a *activityOrderApi) Refund(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceOrder.AdminRefund(idReq.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("退款成功或已受理", c)
}

// GetOrderByCodePublic 公开：根据订单号查询活动订单（H5 核销）
// @Tags LimitedActivityOrder
// @Summary 根据订单号查询活动订单(公开)
// @Param code query string true "订单号"
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /limitedActivity/order/getOrderByCodePublic [get]
func (a *activityOrderApi) GetOrderByCodePublic(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("订单号不能为空", c)
		return
	}
	order, err := serviceOrder.GetByOrderNoPublic(code)
	if err != nil {
		response.FailWithMessage("订单不存在或已失效", c)
		return
	}
	totalUse := order.TotalUseTimes
	if totalUse <= 0 {
		totalUse = 1
	}
	remaining := totalUse - order.VerifiedTimes
	if remaining < 0 {
		remaining = 0
	}
	verifyRecords, _ := serviceOrder.GetVerifyRecords(order.ID)
	response.OkWithData(gin.H{
		"order":          order,
		"remainingTimes": remaining,
		"verifyRecords":  verifyRecords,
	}, c)
}

// VerifyOrderByCodePublic 公开：根据订单号核销活动订单（H5）
// @Tags LimitedActivityOrder
// @Summary 根据订单号核销活动订单(公开)
// @Param code query string true "订单号"
// @Success 200 {object} response.Response{data=object,msg=string} "核销成功"
// @Router /limitedActivity/order/verifyOrderByCodePublic [post]
func (a *activityOrderApi) VerifyOrderByCodePublic(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("订单号不能为空", c)
		return
	}
	if err := serviceOrder.VerifyOrderByOrderNoPublic(code); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, _ := serviceOrder.GetByOrderNoPublic(code)
	totalUse := order.TotalUseTimes
	if totalUse <= 0 {
		totalUse = 1
	}
	remaining := totalUse - order.VerifiedTimes
	if remaining < 0 {
		remaining = 0
	}
	verifyRecords, _ := serviceOrder.GetVerifyRecords(order.ID)
	response.OkWithData(gin.H{
		"order":          order,
		"remainingTimes": remaining,
		"verifyRecords":  verifyRecords,
	}, c)
}
