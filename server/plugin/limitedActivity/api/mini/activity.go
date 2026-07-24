package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	laRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
	"github.com/gin-gonic/gin"
)

var Activity = new(miniActivityApi)

type miniActivityApi struct{}

// List 小程序-活动列表（仅显示中）
// @Tags        小程序-限时活动
// @Summary     活动列表
// @Description 返回显示中的限时活动，含剩余名额与是否可报名
// @Accept      json
// @Produce     json
// @Param       page query int false "页码"
// @Param       pageSize query int false "每页条数"
// @Param       name query string false "活动名称"
// @Success     200 {object} response.Response{data=response.PageResult,msg=string}
// @Router      /limitedActivity/mini/activity/list [get]
func (a *miniActivityApi) List(c *gin.Context) {
	var req laRequest.ActivitySearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}
	list, total, err := svcActivity.GetMiniList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Detail 小程序-活动详情
// @Tags        小程序-限时活动
// @Summary     活动详情
// @Description 仅返回显示中的活动详情
// @Accept      json
// @Produce     json
// @Param       id query int true "活动ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /limitedActivity/mini/activity/detail [get]
func (a *miniActivityApi) Detail(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	act, err := svcActivity.GetMiniDetail(idReq.ID)
	if err != nil {
		response.FailWithMessage("活动不存在或未开放", c)
		return
	}
	response.OkWithData(act, c)
}
