package mini

import (
	"github.com/gin-gonic/gin"
)

type MiniRouter struct{}

// Init 小程序端路由，挂在 public 组，无需后台 JWT 鉴权
func (r *MiniRouter) Init(public, private *gin.RouterGroup) {
	g := public.Group("mini")
	// 通用
	g.GET("ping", commonApi.Ping)
	// 后续小程序接口在此追加，例如：
	// g.GET("config", commonApi.Config)
	// g.POST("user/login", userApi.Login)
}
