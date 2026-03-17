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
	// 预约
	g.POST("reservation/create", mini.Reservation.Create)
	g.POST("reservation/update", mini.Reservation.Update)
	g.GET("reservation/myList", mini.Reservation.MyList)
	g.GET("reservation/myDetail", mini.Reservation.MyDetail)
	g.POST("reservation/cancel", mini.Reservation.Cancel)
	g.GET("reservation/cancelRule", mini.Reservation.CancelRule)
}
