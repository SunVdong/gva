package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/api/mini"
	"github.com/gin-gonic/gin"
)

type miniRouter struct{}

func (r *miniRouter) Init(public, private *gin.RouterGroup) {
	g := public.Group("limitedActivity").Group("mini").Use(middleware.OptionalJWTAuth())
	g.GET("activity/list", mini.Activity.List)
	g.GET("activity/detail", mini.Activity.Detail)

	auth := public.Group("limitedActivity").Group("mini").Use(middleware.JWTAuth())
	auth.POST("order/create", mini.Order.Create)
	auth.GET("order/myList", mini.Order.MyList)
	auth.GET("order/detail", mini.Order.Detail)
}
