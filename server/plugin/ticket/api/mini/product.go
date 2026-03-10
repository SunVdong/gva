package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Product = new(miniProductApi)

type miniProductApi struct{}

// ListByScenic 小程序-按景区获取门票商品列表（仅启用）
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
