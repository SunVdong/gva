package mini

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	MiniRouter
}

var (
	commonApi = api.ApiGroupApp.MiniApiGroup.CommonApi
	authApi   = api.ApiGroupApp.MiniApiGroup.AuthApi
	payApi    = api.ApiGroupApp.MiniApiGroup.PayApi
)

// Init 注册小程序路由，供 initialize/router 调用
func (r *RouterGroup) Init(public, private *gin.RouterGroup) {
	r.MiniRouter.Init(public, private)
}
