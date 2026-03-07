package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

var TimeSlot = new(timeSlotApi)

type timeSlotApi struct{}

// CreateTimeSlot 创建时段
// @Tags CampingTimeSlot
// @Summary 创建预约时段
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CampingTimeSlot true "时段信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /camping/timeSlot/createTimeSlot [post]
func (a *timeSlotApi) CreateTimeSlot(c *gin.Context) {
	var m model.CampingTimeSlot
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceSlot.CreateTimeSlot(&m); err != nil {
		response.FailWithMessage("创建失败", c)
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
	if err := serviceSlot.DeleteTimeSlot(idReq.ID); err != nil {
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
	if err := serviceSlot.DeleteTimeSlotByIds(ids); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateTimeSlot 更新时段
// @Tags CampingTimeSlot
// @Summary 更新时段
// @Security ApiKeyAuth
// @Param data body model.CampingTimeSlot true "时段信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /camping/timeSlot/updateTimeSlot [put]
func (a *timeSlotApi) UpdateTimeSlot(c *gin.Context) {
	var m model.CampingTimeSlot
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceSlot.UpdateTimeSlot(m); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindTimeSlot 根据ID查询时段
// @Tags CampingTimeSlot
// @Summary 根据ID查询时段
// @Security ApiKeyAuth
// @Param id query int true "时段ID"
// @Success 200 {object} response.Response{data=model.CampingTimeSlot,msg=string} "查询成功"
// @Router /camping/timeSlot/findTimeSlot [get]
func (a *timeSlotApi) FindTimeSlot(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceSlot.GetTimeSlot(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(res, c)
}

// GetTimeSlotList 分页获取时段列表
// @Tags CampingTimeSlot
// @Summary 分页获取时段列表
// @Security ApiKeyAuth
// @Param data query request.CampingTimeSlotSearch true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /camping/timeSlot/getTimeSlotList [get]
func (a *timeSlotApi) GetTimeSlotList(c *gin.Context) {
	var req campingRequest.CampingTimeSlotSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceSlot.GetTimeSlotList(req)
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

// GetAllTimeSlotsPublic 公开：获取全部时段
// @Tags CampingTimeSlot
// @Summary 获取全部时段(公开)
// @Success 200 {object} response.Response{data=[]model.CampingTimeSlot,msg=string} "获取成功"
// @Router /camping/timeSlot/getAllTimeSlotsPublic [get]
func (a *timeSlotApi) GetAllTimeSlotsPublic(c *gin.Context) {
	list, err := serviceSlot.GetAllTimeSlots()
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
