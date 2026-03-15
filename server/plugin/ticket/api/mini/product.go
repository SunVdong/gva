package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Product = new(miniProductApi)

type miniProductApi struct{}

// ListByScenic 小程序-按景区获取门票商品列表（仅启用）
// @Tags        小程序
// @Summary     门票商品列表
// @Description 小程序端按景区获取已启用的门票商品列表，分页
// @Accept      json
// @Produce     json
// @Param       scenicId query int  true  "景区ID"
// @Param       page     query int  false "页码"
// @Param       pageSize query int  false "每页条数"
// @Success     200      {object} response.Response{data=response.PageResult,msg=string}
// @Router      /ticket/mini/product/list [get]
func (a *miniProductApi) ListByScenic(c *gin.Context) {
	var req struct {
		ScenicID uint `form:"scenicId" binding:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("请选择景区", c)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}
	status := 1
	searchReq := request.TicketProductSearch{
		ScenicID: req.ScenicID,
		Status:   &status,
	}
	searchReq.Page = req.Page
	searchReq.PageSize = req.PageSize
	list, total, err := svcProduct.GetList(searchReq)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Detail 小程序-商品详情（含 SKU 列表、规则，仅启用）
// @Tags        小程序
// @Summary     门票商品详情
// @Description 小程序端获取商品详情，含 SKU 与规则，仅启用时返回
// @Accept      json
// @Produce     json
// @Param       id query int true "商品ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /ticket/mini/product/detail [get]
func (a *miniProductApi) Detail(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	product, err := svcProduct.Get(idReq.ID)
	if err != nil {
		response.FailWithMessage("商品不存在", c)
		return
	}
	if product.Status != 1 {
		response.FailWithMessage("商品已下架", c)
		return
	}
	status := 1
	skuList, _, _ := svcSku.GetList(request.TicketSkuSearch{ProductID: idReq.ID, Status: &status})
	rules, _ := svcRule.GetByProduct(idReq.ID)
	response.OkWithData(gin.H{
		"product": product,
		"skus":    skuList,
		"rules":   rules,
	}, c)
}
