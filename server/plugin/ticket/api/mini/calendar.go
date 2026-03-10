package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Calendar = new(miniCalendarApi)

type miniCalendarApi struct{}

// GetBySku 小程序-查询某 SKU 的日历可售情况（用于选择游玩日期）
func (a *miniCalendarApi) GetBySku(c *gin.Context) {
	var req request.TicketCalendarSearch
	if err := c.ShouldBindQuery(&req); err != nil || req.SkuID == 0 {
		response.FailWithMessage("请选择门票", c)
		return
	}
	if req.StartDate == "" || req.EndDate == "" {
		response.FailWithMessage("请选择日期范围", c)
		return
	}
	req.Page = 1
	req.PageSize = 90
	list, _, err := svcCalendar.GetBySku(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
