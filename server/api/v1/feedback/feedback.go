package feedback

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	feedbackModel "github.com/flipped-aurora/gin-vue-admin/server/model/feedback"
	feedbackReq "github.com/flipped-aurora/gin-vue-admin/server/model/feedback/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FeedbackApi struct{}

func (a *FeedbackApi) DeleteFeedback(c *gin.Context) {
	var row feedbackModel.Feedback
	if err := c.ShouldBindJSON(&row); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(row.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := feedbackService.DeleteFeedback(row); err != nil {
		global.GVA_LOG.Error("删除失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *FeedbackApi) DeleteFeedbackByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := feedbackService.DeleteFeedbackByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除失败", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

func (a *FeedbackApi) FindFeedback(c *gin.Context) {
	var q feedbackModel.Feedback
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(q.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := feedbackService.GetFeedback(q.ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(row, "查询成功", c)
}

func (a *FeedbackApi) GetFeedbackList(c *gin.Context) {
	var pageInfo feedbackReq.FeedbackSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := feedbackService.GetFeedbackList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败", zap.Error(err))
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
