package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Guide = new(guide)

type guide struct{}

// CreateGuide 创建活动指南
func (a *guide) CreateGuide(c *gin.Context) {
	var guide model.ActivityGuide
	if err := c.ShouldBindJSON(&guide); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := model.ValidateActivityGuideCoverImage(guide.CoverImage); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := model.ValidateActivityGuideMedia(guide.Media); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceGuide.CreateGuide(&guide); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteGuide 删除活动指南
func (a *guide) DeleteGuide(c *gin.Context) {
	ID := c.Query("ID")
	if err := serviceGuide.DeleteGuide(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteGuideByIds 批量删除活动指南
func (a *guide) DeleteGuideByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := serviceGuide.DeleteGuideByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateGuide 更新活动指南
func (a *guide) UpdateGuide(c *gin.Context) {
	var guide model.ActivityGuide
	if err := c.ShouldBindJSON(&guide); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := model.ValidateActivityGuideCoverImage(guide.CoverImage); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := model.ValidateActivityGuideMedia(guide.Media); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceGuide.UpdateGuide(guide); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindGuide 根据ID获取活动指南
func (a *guide) FindGuide(c *gin.Context) {
	ID := c.Query("ID")
	data, err := serviceGuide.GetGuide(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetGuideList 分页获取活动指南列表
func (a *guide) GetGuideList(c *gin.Context) {
	var pageInfo request.GuideSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceGuide.GetGuideList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
