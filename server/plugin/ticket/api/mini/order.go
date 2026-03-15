package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Order = new(miniOrderApi)

type miniOrderApi struct{}

// Create 小程序-提交订单
// @Tags        小程序-景点
// @Summary     提交订单
// @Description 小程序端提交门票订单
// @Accept      json
// @Produce     json
// @Param       data body request.MiniOrderCreate true "订单信息"
// @Success     200  {object} response.Response{data=object,msg=string}
// @Router      /ticket/mini/order/create [post]
func (a *miniOrderApi) Create(c *gin.Context) {
	var req request.MiniOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, err := svcOrder.CreateOrder(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(order, c)
}

// MyList 小程序-我的订单列表
// @Tags        小程序-景点
// @Summary     我的订单列表
// @Description 小程序端获取当前用户的订单列表，分页
// @Accept      json
// @Produce     json
// @Param       userId   query int false "用户ID（需登录）"
// @Param       page     query int false "页码"
// @Param       pageSize query int false "每页条数"
// @Success     200      {object} response.Response{data=response.PageResult,msg=string}
// @Router      /ticket/mini/order/myList [get]
func (a *miniOrderApi) MyList(c *gin.Context) {
	var req request.TicketOrderSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.UserID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}
	list, total, err := svcOrder.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Detail 小程序-订单详情（含订单项）
// @Tags        小程序-景点
// @Summary     订单详情
// @Description 小程序端获取订单详情及订单项
// @Accept      json
// @Produce     json
// @Param       id query int true "订单ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /ticket/mini/order/detail [get]
func (a *miniOrderApi) Detail(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order, items, err := svcOrder.GetByID(idReq.ID)
	if err != nil {
		response.FailWithMessage("订单不存在", c)
		return
	}
	response.OkWithData(gin.H{
		"order": order,
		"items": items,
	}, c)
}
