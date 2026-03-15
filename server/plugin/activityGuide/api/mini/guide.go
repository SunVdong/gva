package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/model/request"
	"github.com/gin-gonic/gin"
	"strconv"
)

type guideApi struct{}

// List 小程序-活动指南列表（仅显示已上架的）
// @Tags        小程序
// @Summary     活动指南列表
// @Description 小程序端获取已上架的活动指南列表，分页
// @Accept      json
// @Produce     json
// @Param       page     query int false "页码"
// @Param       pageSize query int false "每页条数"
// @Success     200      {object} response.Response{data=response.PageResult,msg=string}
// @Router      /activityGuide/mini/guide/list [get]
func (a *guideApi) List(c *gin.Context) {
	var req request.GuideSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	showTrue := true
	req.ShowStatus = &showTrue
	list, total, err := svcGuide.GetGuideList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Detail 小程序-活动指南详情（仅已上架时返回）
// @Tags        小程序
// @Summary     活动指南详情
// @Description 小程序端获取活动指南详情，仅已上架时返回
// @Accept      json
// @Produce     json
// @Param       id query int true "活动指南ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /activityGuide/mini/guide/detail [get]
func (a *guideApi) Detail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("请传入活动指南 id", c)
		return
	}
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("id 格式错误", c)
		return
	}
	data, err := svcGuide.GetGuide(idStr)
	if err != nil {
		response.FailWithMessage("活动指南不存在", c)
		return
	}
	if data.ShowStatus != nil && !*data.ShowStatus {
		response.FailWithMessage("活动指南已下架", c)
		return
	}
	response.OkWithData(data, c)
}
