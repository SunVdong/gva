package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

var VenueOpenTime = new(venueOpenTimeApi)

type venueOpenTimeApi struct{}

// GetByVenue 根据场地ID获取开放时间列表
// @Tags VenueOpenTime
// @Summary 获取场地开放时间
// @Security ApiKeyAuth
// @Param venueId query int true "场地ID"
// @Success 200 {object} response.Response{data=[]model.VenueOpenTime,msg=string} "获取成功"
// @Router /camping/venueOpenTime/getByVenue [get]
func (a *venueOpenTimeApi) GetByVenue(c *gin.Context) {
	var idReq struct {
		VenueID uint `form:"venueId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("场地ID不能为空", c)
		return
	}
	list, err := serviceVenueOpenTime.GetVenueOpenTimeListByVenue(idReq.VenueID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// Save 保存场地的开放时间（覆盖）
// @Tags VenueOpenTime
// @Summary 保存场地开放时间
// @Security ApiKeyAuth
// @Param data body object true "body: { venueId, list: [{ weekDay, openTime, closeTime }] }"
// @Success 200 {object} response.Response{msg=string} "保存成功"
// @Router /camping/venueOpenTime/save [post]
func (a *venueOpenTimeApi) Save(c *gin.Context) {
	var body struct {
		VenueID uint `json:"venueId" binding:"required"`
		List    []struct {
			WeekDay   int    `json:"weekDay"`
			OpenTime  string `json:"openTime"`
			CloseTime string `json:"closeTime"`
		} `json:"list"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list := make([]campingRequest.VenueOpenTimeBody, 0, len(body.List))
	for _, item := range body.List {
		list = append(list, campingRequest.VenueOpenTimeBody{
			VenueID:   body.VenueID,
			WeekDay:   item.WeekDay,
			OpenTime:  item.OpenTime,
			CloseTime: item.CloseTime,
		})
	}
	if err := serviceVenueOpenTime.SaveVenueOpenTimes(body.VenueID, list); err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}
