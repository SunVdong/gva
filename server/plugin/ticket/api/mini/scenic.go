package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Scenic = new(miniScenicApi)

type miniScenicApi struct{}

// List 小程序-景区列表（仅启用）
func (a *miniScenicApi) List(c *gin.Context) {
	status := 1
	req := request.ScenicSearch{Status: &status}
	req.Page = 1
	req.PageSize = 20
	_ = c.ShouldBindQuery(&req)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	list, total, err := svcScenic.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, "获取成功", c)
}

// Detail 小程序-景区详情（仅启用时返回）
func (a *miniScenicApi) Detail(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	res, err := svcScenic.Get(idReq.ID)
	if err != nil {
		response.FailWithMessage("景区不存在", c)
		return
	}
	if res.Status != 1 {
		response.FailWithMessage("景区已下架", c)
		return
	}
	openTimes, _ := svcOpenTime.GetByScenic(idReq.ID)
	response.OkWithData(gin.H{
		"scenic":    res,
		"openTimes": openTimes,
	}, c)
}
