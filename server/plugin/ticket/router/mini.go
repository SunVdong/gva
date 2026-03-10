package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/api/mini"
	"github.com/gin-gonic/gin"
)

type miniRouter struct{}

// Init 小程序端路由，全部挂在 public，无需后台 JWT
func (r *miniRouter) Init(public, private *gin.RouterGroup) {
	g := public.Group("ticket").Group("mini")
	// 景区
	g.GET("scenic/list", mini.Scenic.List)
	g.GET("scenic/detail", mini.Scenic.Detail)
	// 门票商品
	g.GET("product/list", mini.Product.ListByScenic)
	g.GET("product/detail", mini.Product.Detail)
	// 日历可售
	g.GET("calendar/sku", mini.Calendar.GetBySku)
	// 订单
	g.POST("order/create", mini.Order.Create)
	g.GET("order/myList", mini.Order.MyList)
	g.GET("order/detail", mini.Order.Detail)
}
