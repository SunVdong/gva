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
	order, items, err := serviceOrder.GetByID(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(gin.H{
		"order": order,
		"items": items,
	}, c)
}
