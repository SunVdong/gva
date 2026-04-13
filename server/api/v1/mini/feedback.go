package mini

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	feedbackModel "github.com/flipped-aurora/gin-vue-admin/server/model/feedback"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type FeedbackApi struct{}

// SubmitFeedback 提交意见反馈
// @Tags        小程序
// @Summary     新增意见反馈
// @Description 须先登录，请求头携带 x-token；正文 1～2000 字
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token（必填，与 wx 请求头一致）"
// @Param       data body object true "请求体" example({"content":"反馈正文"})
// @Success     200 {object} response.Response{msg=string} "提交成功"
// @Router      /mini/feedback [post]
func (a *FeedbackApi) SubmitFeedback(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请填写反馈内容", c)
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if len(req.Content) == 0 || len(req.Content) > 2000 {
		response.FailWithMessage("反馈内容长度需在 1～2000 字", c)
		return
	}
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithMessage("请先登录后再反馈", c)
		return
	}
	f := feedbackModel.Feedback{
		UserID:  uid,
		Content: req.Content,
	}
	if err := service.ServiceGroupApp.FeedbackServiceGroup.FeedbackService.CreateFeedback(f); err != nil {
		response.FailWithMessage("提交失败", c)
		return
	}
	response.OkWithMessage("感谢您的反馈", c)
}

// ListMyFeedback 查看本人历史意见反馈
// @Tags        小程序
// @Summary     本人反馈列表
// @Description 分页返回当前登录用户的反馈记录，仅本人；须请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token（必填）"
// @Param       page query int false "页码，默认 1"
// @Param       pageSize query int false "每页条数，默认 10，最大 100"
// @Success     200 {object} response.Response{data=response.PageResult,msg=string} "data.list 为反馈列表项（含 ID、content、CreatedAt 等）"
// @Router      /mini/feedback [get]
func (a *FeedbackApi) ListMyFeedback(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	list, total, err := service.ServiceGroupApp.FeedbackServiceGroup.FeedbackService.ListFeedbackByUserID(uid, pageInfo)
	if err != nil {
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
