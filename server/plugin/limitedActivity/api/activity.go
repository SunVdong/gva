package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model"
	laRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var Activity = new(activityApi)

type activityApi struct{}

// Create 创建限时活动
// @Tags LimitedActivity
// @Summary 创建限时活动
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Activity true "活动信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /limitedActivity/activity/createActivity [post]
func (a *activityApi) Create(c *gin.Context) {
	var m model.Activity
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	m.CreatedBy = int(utils.GetUserID(c))
	m.UpdatedBy = m.CreatedBy
	if err := serviceActivity.Create(&m); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// Delete 删除限时活动
// @Tags LimitedActivity
// @Summary 删除限时活动
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query int true "活动ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /limitedActivity/activity/deleteActivity [delete]
func (a *activityApi) Delete(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceActivity.Delete(idReq.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteByIds 批量删除限时活动
// @Tags LimitedActivity
// @Summary 批量删除限时活动
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body []uint true "ID列表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /limitedActivity/activity/deleteActivityByIds [delete]
func (a *activityApi) DeleteByIds(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceActivity.DeleteByIds(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// Update 更新限时活动
// @Tags LimitedActivity
// @Summary 更新限时活动
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Activity true "活动信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /limitedActivity/activity/updateActivity [put]
func (a *activityApi) Update(c *gin.Context) {
	var m model.Activity
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	m.UpdatedBy = int(utils.GetUserID(c))
	if err := serviceActivity.Update(m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// Find 查询限时活动
// @Tags LimitedActivity
// @Summary 查询限时活动
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query int true "活动ID"
// @Success 200 {object} response.Response{data=model.Activity,msg=string} "查询成功"
// @Router /limitedActivity/activity/findActivity [get]
func (a *activityApi) Find(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceActivity.Get(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(res, c)
}

// GetList 限时活动列表
// @Tags LimitedActivity
// @Summary 限时活动列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param data query laRequest.ActivitySearch true "分页搜索"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /limitedActivity/activity/getActivityList [get]
func (a *activityApi) GetList(c *gin.Context) {
	var req laRequest.ActivitySearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceActivity.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}
