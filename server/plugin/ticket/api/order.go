package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	ticketRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Order = new(ticketOrderApi)

type ticketOrderApi struct{}

func (a *ticketOrderApi) GetList(c *gin.Context) {
	var req ticketRequest.TicketOrderSearch
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
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

func (a *ticketOrderApi) Find(c *gin.Context) {
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
	data := gin.H{
		"order": order,
	}
	if order.Status == 2 && order.VerifiedAt != nil {
		review, _ := serviceOrderReview.GetByOrderID(order.ID)
		if review.ID != 0 {
			data["review"] = gin.H{
				"id":        review.ID,
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

// GetOrderByCodePublic 公开：根据订单号(code)查询门票订单（用于 H5 扫码核销）
// @Tags TicketOrder
// @Summary 根据订单号查询门票订单(公开)
// @Param code query string true "订单号"
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /ticket/order/getOrderByCodePublic [get]
func (a *ticketOrderApi) GetOrderByCodePublic(c *gin.Context) {
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
	response.OkWithData(gin.H{
		"order": order,
	}, c)
}

// VerifyOrderByCodePublic 公开：根据订单号(code)核销门票订单（用于 H5 扫码核销）
// @Tags TicketOrder
// @Summary 根据订单号核销门票订单(公开)
// @Param code query string true "订单号"
// @Success 200 {object} response.Response{data=object,msg=string} "核销成功"
// @Router /ticket/order/verifyOrderByCodePublic [post]
func (a *ticketOrderApi) VerifyOrderByCodePublic(c *gin.Context) {
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
	response.OkWithData(gin.H{
		"order": order,
	}, c)
}
