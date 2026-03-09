package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

var TimeSlot = new(timeSlotApi)

type timeSlotApi struct{}

// CreateTimeSlot 创建场地时间段
// @Tags CampingTimeSlot
// @Summary 创建预约时段（需指定场地ID）
// @Security ApiKeyAuth
// @Param data body model.VenueTimeslot true "时段信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /camping/timeSlot/createTimeSlot [post]
func (a *timeSlotApi) CreateTimeSlot(c *gin.Context) {
	var m model.VenueTimeslot
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenueTimeslot.CreateVenueTimeslot(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteTimeSlot 删除时段
// @Tags CampingTimeSlot
// @Summary 删除时段
// @Security ApiKeyAuth
// @Param id query int true "时段ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /camping/timeSlot/deleteTimeSlot [delete]
func (a *timeSlotApi) DeleteTimeSlot(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenueTimeslot.DeleteVenueTimeslot(idReq.ID); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteTimeSlotByIds 批量删除时段
// @Tags CampingTimeSlot
// @Summary 批量删除时段
// @Security ApiKeyAuth
// @Param data body []uint true "ID数组"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /camping/timeSlot/deleteTimeSlotByIds [delete]
func (a *timeSlotApi) DeleteTimeSlotByIds(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenueTimeslot.DeleteVenueTimeslotByIds(ids); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateTimeSlot 更新时段
// @Tags CampingTimeSlot
// @Summary 更新时段
// @Security ApiKeyAuth
// @Param data body model.VenueTimeslot true "时段信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /camping/timeSlot/updateTimeSlot [put]
func (a *timeSlotApi) UpdateTimeSlot(c *gin.Context) {
	var m model.VenueTimeslot
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenueTimeslot.UpdateVenueTimeslot(m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindTimeSlot 根据ID查询时段
// @Tags CampingTimeSlot
// @Summary 根据ID查询时段
// @Security ApiKeyAuth
// @Param id query int true "时段ID"
// @Success 200 {object} response.Response{data=model.VenueTimeslot,msg=string} "查询成功"
// @Router /camping/timeSlot/findTimeSlot [get]
func (a *timeSlotApi) FindTimeSlot(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceVenueTimeslot.GetVenueTimeslot(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(res, c)
}

// GetTimeSlotList 分页获取时段列表（可按场地ID筛选）
// @Tags CampingTimeSlot
// @Summary 分页获取时段列表
// @Security ApiKeyAuth
// @Param data query request.VenueTimeslotSearch true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /camping/timeSlot/getTimeSlotList [get]
func (a *timeSlotApi) GetTimeSlotList(c *gin.Context) {
	var req campingRequest.VenueTimeslotSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceVenueTimeslot.GetVenueTimeslotList(req)
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

// GetAllTimeSlotsPublic 公开：获取某场地的全部时间段（兼容旧接口，无 venueId 时返回空）
// @Tags CampingTimeSlot
// @Summary 获取场地时间段列表(公开)
// @Param venueId query int false "场地ID，不传则返回空"
// @Success 200 {object} response.Response{data=[]model.VenueTimeslot,msg=string} "获取成功"
// @Router /camping/timeSlot/getAllTimeSlotsPublic [get]
func (a *timeSlotApi) GetAllTimeSlotsPublic(c *gin.Context) {
	venueIDStr := c.Query("venueId")
	if venueIDStr == "" {
		response.OkWithData([]model.VenueTimeslot{}, c)
		return
	}
	var idReq struct {
		VenueID uint `form:"venueId"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil || idReq.VenueID == 0 {
		response.OkWithData([]model.VenueTimeslot{}, c)
		return
	}
	list, err := serviceVenueTimeslot.GetVenueTimeslotsByVenue(idReq.VenueID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// GetTimeSlotsByVenuePublic 公开：获取某场地的全部时间段
// @Tags CampingTimeSlot
// @Summary 获取某场地时间段列表(公开)
// @Param venueId query int true "场地ID"
// @Success 200 {object} response.Response{data=[]model.VenueTimeslot,msg=string} "获取成功"
// @Router /camping/timeSlot/getTimeSlotsByVenuePublic [get]
func (a *timeSlotApi) GetTimeSlotsByVenuePublic(c *gin.Context) {
	var idReq struct {
		VenueID uint `form:"venueId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("场地ID不能为空", c)
		return
	}
	list, err := serviceVenueTimeslot.GetVenueTimeslotsByVenue(idReq.VenueID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
