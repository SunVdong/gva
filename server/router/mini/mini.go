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
	// 登录（微信 code2session + JWT）
	g.POST("login", authApi.Login)
}
