package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/api"
	"github.com/gin-gonic/gin"
)

var apiBanner = api.Api.Banner

type bannerRouter struct{}

func (r *bannerRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("limitedActivity").Group("banner")
	g.Use(middleware.OperationRecord()).POST("createBanner", apiBanner.Create)
	g.Use(middleware.OperationRecord()).DELETE("deleteBanner", apiBanner.Delete)
	g.Use(middleware.OperationRecord()).DELETE("deleteBannerByIds", apiBanner.DeleteByIds)
	g.Use(middleware.OperationRecord()).PUT("updateBanner", apiBanner.Update)
	g.GET("findBanner", apiBanner.Find)
	g.GET("getBannerList", apiBanner.GetList)
}
