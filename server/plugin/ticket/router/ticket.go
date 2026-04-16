package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/api"
	"github.com/gin-gonic/gin"
)

var (
	apiScenic         = api.Api.Scenic
	apiScenicOpenTime = api.Api.ScenicOpenTime
	apiProduct        = api.Api.Product
	apiSku            = api.Api.Sku
	apiRule           = api.Api.Rule
	apiCalendar       = api.Api.Calendar
	apiUser           = api.Api.User
	apiOrder          = api.Api.Order
)

type scenicRouter struct{}
type scenicOpenTimeRouter struct{}
type productRouter struct{}
type skuRouter struct{}
type ruleRouter struct{}
type calendarRouter struct{}
type userRouter struct{}
type orderRouter struct{}

func (r *scenicRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("scenic")
	g.Use(middleware.OperationRecord()).POST("createScenic", apiScenic.Create)
	g.Use(middleware.OperationRecord()).DELETE("deleteScenic", apiScenic.Delete)
	g.Use(middleware.OperationRecord()).DELETE("deleteScenicByIds", apiScenic.DeleteByIds)
	g.Use(middleware.OperationRecord()).PUT("updateScenic", apiScenic.Update)
	g.GET("findScenic", apiScenic.Find)
	g.GET("getScenicList", apiScenic.GetList)
}

func (r *scenicOpenTimeRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("scenicOpenTime")
	g.GET("getByScenic", apiScenicOpenTime.GetByScenic)
	g.Use(middleware.OperationRecord()).POST("save", apiScenicOpenTime.Save)
}

func (r *productRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("product")
	g.Use(middleware.OperationRecord()).POST("createProduct", apiProduct.Create)
	g.Use(middleware.OperationRecord()).DELETE("deleteProduct", apiProduct.Delete)
	g.Use(middleware.OperationRecord()).DELETE("deleteProductByIds", apiProduct.DeleteByIds)
	g.Use(middleware.OperationRecord()).PUT("updateProduct", apiProduct.Update)
	g.GET("findProduct", apiProduct.Find)
	g.GET("getProductList", apiProduct.GetList)
}

func (r *skuRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("sku")
	g.Use(middleware.OperationRecord()).POST("createSku", apiSku.Create)
	g.Use(middleware.OperationRecord()).DELETE("deleteSku", apiSku.Delete)
	g.Use(middleware.OperationRecord()).DELETE("deleteSkuByIds", apiSku.DeleteByIds)
	g.Use(middleware.OperationRecord()).PUT("updateSku", apiSku.Update)
	g.GET("findSku", apiSku.Find)
	g.GET("getSkuList", apiSku.GetList)
}

func (r *ruleRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("rule")
	g.GET("getByProduct", apiRule.GetByProduct)
	g.Use(middleware.OperationRecord()).POST("save", apiRule.Save)
}

func (r *calendarRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("calendar")
	g.GET("getBySku", apiCalendar.GetBySku)
	g.Use(middleware.OperationRecord()).POST("set", apiCalendar.Set)
}

func (r *userRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("user")
	g.GET("getUserList", apiUser.GetList)
	g.GET("findUser", apiUser.Find)
	g.Use(middleware.OperationRecord()).PUT("updateUser", apiUser.Update)
}

func (r *orderRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("ticket").Group("order")
	g.Use(middleware.OperationRecord()).POST("refundOrder", apiOrder.Refund)
	g.GET("getOrderList", apiOrder.GetList)
	g.GET("findOrder", apiOrder.Find)

	// H5 核销公开接口（根据订单号查询与核销）
	pg := public.Group("ticket").Group("order")
	pg.GET("getOrderByCodePublic", apiOrder.GetOrderByCodePublic)
	pg.POST("verifyOrderByCodePublic", apiOrder.VerifyOrderByCodePublic)
}
