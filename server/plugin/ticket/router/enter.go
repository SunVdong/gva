package router

import "github.com/gin-gonic/gin"

var Router = new(router)

type router struct {
	Scenic         scenicRouter
	ScenicOpenTime scenicOpenTimeRouter
	Product        productRouter
	Sku            skuRouter
	Rule           ruleRouter
	Audience       audienceRouter
	Calendar       calendarRouter
	User           userRouter
	Order          orderRouter
	Mini           miniRouter
}

// Init 初始化门票插件路由
func (r *router) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	r.Scenic.Init(public, private)
	r.ScenicOpenTime.Init(public, private)
	r.Product.Init(public, private)
	r.Sku.Init(public, private)
	r.Rule.Init(public, private)
	r.Audience.Init(public, private)
	r.Calendar.Init(public, private)
	r.User.Init(public, private)
	r.Order.Init(public, private)
	r.Mini.Init(public, private)
}
