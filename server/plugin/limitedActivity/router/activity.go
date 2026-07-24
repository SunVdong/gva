package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/api"
	"github.com/gin-gonic/gin"
)

var (
	apiActivity = api.Api.Activity
	apiOrder    = api.Api.Order
)

type activityRouter struct{}
type orderRouter struct{}

func (r *activityRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("limitedActivity").Group("activity")
	g.Use(middleware.OperationRecord()).POST("createActivity", apiActivity.Create)
	g.Use(middleware.OperationRecord()).DELETE("deleteActivity", apiActivity.Delete)
	g.Use(middleware.OperationRecord()).DELETE("deleteActivityByIds", apiActivity.DeleteByIds)
	g.Use(middleware.OperationRecord()).PUT("updateActivity", apiActivity.Update)
	g.GET("findActivity", apiActivity.Find)
	g.GET("getActivityList", apiActivity.GetList)
}

func (r *orderRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("limitedActivity").Group("order")
	g.Use(middleware.OperationRecord()).POST("refundOrder", apiOrder.Refund)
	g.GET("getOrderList", apiOrder.GetList)
	g.GET("findOrder", apiOrder.Find)

	pg := public.Group("limitedActivity").Group("order")
	pg.GET("getOrderByCodePublic", apiOrder.GetOrderByCodePublic)
	pg.POST("verifyOrderByCodePublic", apiOrder.VerifyOrderByCodePublic)
}
