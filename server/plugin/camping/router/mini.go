package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/api/mini"
	"github.com/gin-gonic/gin"
)

type miniRouter struct{}

// Init 小程序端露营预约路由，挂在 public 组
// 需登录的接口（我的预约列表/详情/取消）依赖中间件或网关注入 x-user-id
func (r *miniRouter) Init(public, private *gin.RouterGroup) {
	g := public.Group("camping").Group("mini")
	// 场地
	g.GET("site/list", mini.Site.List)
	g.GET("site/detail", mini.Site.Detail)
	// 时段
	g.GET("slot/availableSlots", mini.Slot.AvailableSlots)
	// 预约
	g.POST("reservation/create", mini.Reservation.Create)
	g.GET("reservation/myList", mini.Reservation.MyList)
	g.GET("reservation/myDetail", mini.Reservation.MyDetail)
	g.POST("reservation/cancel", mini.Reservation.Cancel)
	g.GET("reservation/cancelRule", mini.Reservation.CancelRule)
}
