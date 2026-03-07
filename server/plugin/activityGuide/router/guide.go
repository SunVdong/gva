package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/api"
	"github.com/gin-gonic/gin"
)

var Guide = new(guide)

type guide struct{}

func (r *guide) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("activityGuide").Use(middleware.OperationRecord())
		group.POST("createGuide", api.Guide.CreateGuide)
		group.DELETE("deleteGuide", api.Guide.DeleteGuide)
		group.DELETE("deleteGuideByIds", api.Guide.DeleteGuideByIds)
		group.PUT("updateGuide", api.Guide.UpdateGuide)
	}
	{
		group := private.Group("activityGuide")
		group.GET("findGuide", api.Guide.FindGuide)
		group.GET("getGuideList", api.Guide.GetGuideList)
	}
}
