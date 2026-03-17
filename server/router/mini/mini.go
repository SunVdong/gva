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
	g.GET("genToken", commonApi.GenToken)
	// 登录（组合登录：wx.login code + getPhoneNumber code）
	g.POST("login", authApi.Login)

	// 微信支付（公共接口，景点/露营等均可复用）
	g.POST("pay/create", payApi.Create)   // 调起支付，需登录，返回 wx.requestPayment 参数
	g.POST("pay/notify", payApi.Notify)    // 支付结果回调，由微信服务器调用
}
