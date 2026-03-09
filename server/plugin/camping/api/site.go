package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

var Site = new(siteApi)

type siteApi struct{}

// CreateSite 创建场地（Venue）
// @Tags CampingSite
// @Summary 创建露营场地
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Venue true "场地信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /camping/site/createSite [post]
func (a *siteApi) CreateSite(c *gin.Context) {
	var m model.Venue
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenue.CreateVenue(&m); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteSite 删除场地
// @Tags CampingSite
// @Summary 删除场地
// @Security ApiKeyAuth
// @Param id query int true "场地ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /camping/site/deleteSite [delete]
func (a *siteApi) DeleteSite(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenue.DeleteVenue(idReq.ID); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSiteByIds 批量删除场地
// @Tags CampingSite
// @Summary 批量删除场地
// @Security ApiKeyAuth
// @Param data body []uint true "ID数组"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /camping/site/deleteSiteByIds [delete]
func (a *siteApi) DeleteSiteByIds(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenue.DeleteVenueByIds(ids); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateSite 更新场地
// @Tags CampingSite
// @Summary 更新场地
// @Security ApiKeyAuth
// @Param data body model.Venue true "场地信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /camping/site/updateSite [put]
func (a *siteApi) UpdateSite(c *gin.Context) {
	var m model.Venue
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceVenue.UpdateVenue(m); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSite 根据ID查询场地
// @Tags CampingSite
// @Summary 根据ID查询场地
// @Security ApiKeyAuth
// @Param id query int true "场地ID"
// @Success 200 {object} response.Response{data=model.Venue,msg=string} "查询成功"
// @Router /camping/site/findSite [get]
func (a *siteApi) FindSite(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceVenue.GetVenue(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(res, c)
}

// GetSiteList 分页获取场地列表
// @Tags CampingSite
// @Summary 分页获取场地列表
// @Security ApiKeyAuth
// @Param data query request.VenueSearch true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /camping/site/getSiteList [get]
func (a *siteApi) GetSiteList(c *gin.Context) {
	var req campingRequest.VenueSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceVenue.GetVenueList(req)
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

// GetSiteListPublic 公开接口：获取启用场地列表（用于前台预约页）
// @Tags CampingSite
// @Summary 获取启用场地列表(公开)
// @Success 200 {object} response.Response{data=[]model.Venue,msg=string} "获取成功"
// @Router /camping/site/getSiteListPublic [get]
func (a *siteApi) GetSiteListPublic(c *gin.Context) {
	status := 1
	req := campingRequest.VenueSearch{
		Status:   &status,
		PageInfo: request.PageInfo{Page: 1, PageSize: 100},
	}
	list, _, err := serviceVenue.GetVenueList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// GetSiteDetailPublic 公开接口：根据ID获取场地详情
// @Tags CampingSite
// @Summary 获取场地详情(公开)
// @Param id query int true "场地ID"
// @Success 200 {object} response.Response{data=model.Venue,msg=string} "获取成功"
// @Router /camping/site/getSiteDetailPublic [get]
func (a *siteApi) GetSiteDetailPublic(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceVenue.GetVenue(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	if res.Status != 1 {
		response.FailWithMessage("场地已禁用", c)
		return
	}
	response.OkWithData(res, c)
}
