package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/api/mini"
	"github.com/gin-gonic/gin"
)

type miniRouter struct{}

// Init 小程序端路由，挂在 public 组，无需后台 JWT
func (r *miniRouter) Init(public, private *gin.RouterGroup) {
	g := public.Group("activityGuide").Group("mini")
	g.GET("guide/list", mini.Guide.List)
	g.GET("guide/detail", mini.Guide.Detail)
}
