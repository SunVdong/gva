package api

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

var VenueCalendar = new(venueCalendarApi)

type venueCalendarApi struct{}

// Set 设置某日是否可预约
// @Tags VenueCalendar
// @Summary 设置场地某日状态
// @Security ApiKeyAuth
// @Param data body request.VenueCalendarSet true "venueId, date(2006-01-02), status"
// @Success 200 {object} response.Response{msg=string} "设置成功"
// @Router /camping/venueCalendar/set [post]
func (a *venueCalendarApi) Set(c *gin.Context) {
	var req campingRequest.VenueCalendarSet
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	date, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		response.FailWithMessage("日期格式错误，请使用 2006-01-02", c)
		return
	}
	if err := serviceVenueCalendar.SetOrCreate(req.VenueID, date, req.Status); err != nil {
		response.FailWithMessage("设置失败", c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// GetByVenue 查询场地某段日期的日历
// @Tags VenueCalendar
// @Summary 获取场地日历
// @Security ApiKeyAuth
// @Param venueId query int true "场地ID"
// @Param start query string true "开始日期 2006-01-02"
// @Param end query string true "结束日期 2006-01-02"
// @Success 200 {object} response.Response{data=[]model.VenueCalendar,msg=string} "获取成功"
// @Router /camping/venueCalendar/getByVenue [get]
func (a *venueCalendarApi) GetByVenue(c *gin.Context) {
	var req campingRequest.VenueCalendarQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	start, err := time.ParseInLocation("2006-01-02", req.Start, time.Local)
	if err != nil {
		response.FailWithMessage("开始日期格式错误", c)
		return
	}
	end, err := time.ParseInLocation("2006-01-02", req.End, time.Local)
	if err != nil {
		response.FailWithMessage("结束日期格式错误", c)
		return
	}
	list, err := serviceVenueCalendar.GetByVenueAndDateRange(req.VenueID, start, end)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
