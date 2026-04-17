package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/api/mini"
	"github.com/gin-gonic/gin"
)

type miniRouter struct{}

// Init 小程序端路由，全部挂在 public。
// 依赖 OptionalJWTAuth：带 x-token 时解析并注入 x-user-id，用于订单等接口做登录校验。
func (r *miniRouter) Init(public, private *gin.RouterGroup) {
	g := public.Group("ticket").Group("mini").Use(middleware.OptionalJWTAuth())
	// 景区
	g.GET("scenic/list", mini.Scenic.List)
	g.GET("scenic/detail", mini.Scenic.Detail)
	// 门票商品
	g.GET("product/list", mini.Product.ListByScenic)
	g.GET("product/detail", mini.Product.Detail)
	// 日历可售
	g.GET("calendar/sku", mini.Calendar.GetBySku)

	// 需登录接口：强制 JWT 鉴权，未登录/过期统一返回 401
	auth := public.Group("ticket").Group("mini").Use(middleware.JWTAuth())
	// 订单
	auth.POST("order/create", mini.Order.Create)
	auth.GET("order/myList", mini.Order.MyList)
	auth.GET("order/detail", mini.Order.Detail)
	auth.POST("order/closePending", mini.Order.ClosePending)
	auth.POST("order/delete", mini.Order.Delete)
	// 订单评价（仅核销后可评价、可删除）
	auth.POST("order/review/create", mini.Order.CreateReview)
	auth.POST("order/review/delete", mini.Order.DeleteReview)
}
