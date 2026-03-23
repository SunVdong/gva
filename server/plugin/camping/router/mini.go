package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/api/mini"
	"github.com/gin-gonic/gin"
)

type miniRouter struct{}

// Init 小程序端露营预约路由，挂在 public 组
// 需登录的接口（我的预约列表/详情/取消/提交）依赖 OptionalJWTAuth 注入 x-user-id
func (r *miniRouter) Init(public, private *gin.RouterGroup) {
	// camping/mini 下挂可选 JWT 中间件：
	// - 带 x-token 时自动解析并注入 x-user-id
	// - 不带 token 时不拦截，可继续访问公开接口
	g := public.Group("camping").Group("mini").Use(middleware.OptionalJWTAuth())
	// 场地
	g.GET("site/list", mini.Site.List)
	g.GET("site/detail", mini.Site.Detail)
	// 时段
	g.GET("slot/availableSlots", mini.Slot.AvailableSlots)
	// 公开接口
	g.GET("reservation/cancelRule", mini.Reservation.CancelRule)

	// 需登录接口：强制 JWT 鉴权，未登录/过期统一返回 401
	auth := g.Group("").Use(middleware.JWTAuth())
	// 预约
	auth.POST("reservation/create", mini.Reservation.Create)
	auth.POST("reservation/update", mini.Reservation.Update)
	auth.GET("reservation/myList", mini.Reservation.MyList)
	auth.GET("reservation/myDetail", mini.Reservation.MyDetail)
	auth.POST("reservation/cancel", mini.Reservation.Cancel)
	// 预约评价（仅核销后可评价、可删除）
	auth.POST("reservation/review/create", mini.Reservation.CreateReview)
	auth.POST("reservation/review/delete", mini.Reservation.DeleteReview)
}
