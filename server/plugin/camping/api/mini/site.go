package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

type siteApi struct{}

// List 小程序-露营场地列表（仅启用）
// @Tags        小程序-露营
// @Summary     场地列表
// @Description 小程序端获取可预约的露营场地列表
// @Accept      json
// @Produce     json
// @Success     200 {object} response.Response{data=[]model.Venue,msg=string}
// @Router      /camping/mini/site/list [get]
func (a *siteApi) List(c *gin.Context) {
	status := 1
	req := campingRequest.VenueSearch{Status: &status}
	req.Page = 1
	req.PageSize = 100
	list, _, err := svcVenue.GetVenueList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// Detail 小程序-场地详情（含介绍、开放时间、预约规则）
// @Tags        小程序-露营
// @Summary     场地详情
// @Description 小程序端获取场地详情，含场地介绍、开放时间、预约规则
// @Accept      json
// @Produce     json
// @Param       id query int true "场地ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /camping/mini/site/detail [get]
func (a *siteApi) Detail(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("请传入场地 id", c)
		return
	}
	venue, err := svcVenue.GetVenue(idReq.ID)
	if err != nil {
		response.FailWithMessage("场地不存在", c)
		return
	}
	if venue.Status != 1 {
		response.FailWithMessage("场地已关闭", c)
		return
	}
	openTimes, _ := svcVenueOpenTime.GetVenueOpenTimeListByVenue(idReq.ID)
	response.OkWithData(gin.H{
		"id":            venue.ID,
		"name":          venue.Name,
		"description":   venue.Description,
		"carouselImages": venue.CarouselImages,
		"reserveRules":  venue.ReserveRules,
		"openTimes":     openTimes,
		"refundChangeHours": venue.RefundChangeHours,
	}, c)
}
