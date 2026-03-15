package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/api"
	"github.com/gin-gonic/gin"
)

var (
	Router    = new(router)
	apiGuide  = api.Api.Guide
)

type router struct {
	Guide guide
	Mini  miniRouter
}

// Init 初始化活动指南路由
func (r *router) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	r.Guide.Init(public, private)
	r.Mini.Init(public, private)
}
