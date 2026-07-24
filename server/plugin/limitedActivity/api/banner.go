package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model"
	laRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var Banner = new(bannerApi)

type bannerApi struct{}

// Create 创建 Banner
// @Tags LimitedActivityBanner
// @Summary 创建 Banner
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Banner true "Banner 信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /limitedActivity/banner/createBanner [post]
func (a *bannerApi) Create(c *gin.Context) {
	var m model.Banner
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	m.CreatedBy = int(utils.GetUserID(c))
	m.UpdatedBy = m.CreatedBy
	if err := serviceBanner.Create(&m); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// Delete 删除 Banner
// @Tags LimitedActivityBanner
// @Summary 删除 Banner
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query int true "Banner ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /limitedActivity/banner/deleteBanner [delete]
func (a *bannerApi) Delete(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceBanner.Delete(idReq.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteByIds 批量删除 Banner
// @Tags LimitedActivityBanner
// @Summary 批量删除 Banner
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body []uint true "ID列表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /limitedActivity/banner/deleteBannerByIds [delete]
func (a *bannerApi) DeleteByIds(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceBanner.DeleteByIds(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// Update 更新 Banner
// @Tags LimitedActivityBanner
// @Summary 更新 Banner
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Banner true "Banner 信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /limitedActivity/banner/updateBanner [put]
func (a *bannerApi) Update(c *gin.Context) {
	var m model.Banner
	if err := c.ShouldBindJSON(&m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	m.UpdatedBy = int(utils.GetUserID(c))
	if err := serviceBanner.Update(m); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// Find 查询 Banner
// @Tags LimitedActivityBanner
// @Summary 查询 Banner
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query int true "Banner ID"
// @Success 200 {object} response.Response{data=model.Banner,msg=string} "查询成功"
// @Router /limitedActivity/banner/findBanner [get]
func (a *bannerApi) Find(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceBanner.Get(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(res, c)
}

// GetList Banner 列表
// @Tags LimitedActivityBanner
// @Summary Banner 列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param data query laRequest.BannerSearch true "分页搜索"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /limitedActivity/banner/getBannerList [get]
func (a *bannerApi) GetList(c *gin.Context) {
	var req laRequest.BannerSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceBanner.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}
