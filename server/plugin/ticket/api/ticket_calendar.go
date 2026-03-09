package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	ticketRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Calendar = new(ticketCalendarApi)

type ticketCalendarApi struct{}

func (a *ticketCalendarApi) GetBySku(c *gin.Context) {
	var req ticketRequest.TicketCalendarSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.SkuID == 0 {
		response.FailWithMessage("skuId不能为空", c)
		return
	}
	list, total, err := serviceCalendar.GetBySku(req)
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

func (a *ticketCalendarApi) Set(c *gin.Context) {
	var req ticketRequest.TicketCalendarSet
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceCalendar.Set(req); err != nil {
		response.FailWithMessage("设置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}
