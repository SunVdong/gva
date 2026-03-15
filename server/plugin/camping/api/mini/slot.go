package mini

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type slotApi struct{}

// AvailableSlots 小程序-某日可预约时段
// @Tags        小程序-露营
// @Summary     某日可预约时段
// @Description 按场地与日期返回时段列表，并标记是否可约（未满且当日开放）
// @Accept      json
// @Produce     json
// @Param       venueId    query int    true "场地ID"
// @Param       reserveDate query string true "预约日期 2006-01-02"
// @Success     200 {object} response.Response{data=[]object,msg=string}
// @Router      /camping/mini/slot/availableSlots [get]
func (a *slotApi) AvailableSlots(c *gin.Context) {
	venueIDStr := c.Query("venueId")
	reserveDateStr := c.Query("reserveDate")
	if venueIDStr == "" || reserveDateStr == "" {
		response.FailWithMessage("场地ID和预约日期不能为空", c)
		return
	}
	venueID, err := strconv.ParseUint(venueIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("场地ID格式错误", c)
		return
	}
	reserveDate, err := time.ParseInLocation("2006-01-02", reserveDateStr, time.Local)
	if err != nil {
		response.FailWithMessage("日期格式错误，请使用 2006-01-02", c)
		return
	}
	open, err := svcVenueCalendar.IsDateOpen(uint(venueID), reserveDate)
	if err != nil || !open {
		response.OkWithData([]gin.H{}, c)
		return
	}
	slots, err := svcVenueTimeslot.GetVenueTimeslotsByVenue(uint(venueID))
	if err != nil {
		response.FailWithMessage("获取时段失败", c)
		return
	}
	fullIDs, err := svcReservation.GetReservedTimeslotIds(uint(venueID), reserveDate)
	if err != nil {
		response.FailWithMessage("获取可约状态失败", c)
		return
	}
	fullMap := make(map[uint]bool)
	for _, id := range fullIDs {
		fullMap[id] = true
	}
	list := make([]gin.H, 0, len(slots))
	for _, s := range slots {
		list = append(list, gin.H{
			"id":         s.ID,
			"startTime":  s.StartTime.FormatHHMM(),
			"endTime":    s.EndTime.FormatHHMM(),
			"capacity":   s.Capacity,
			"available":  !fullMap[s.ID],
		})
	}
	response.OkWithData(list, c)
}
