package router

import "github.com/gin-gonic/gin"

var Router = new(router)

type router struct {
	Activity activityRouter
	Order    orderRouter
	Mini     miniRouter
}

func (r *router) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	r.Activity.Init(public, private)
	r.Order.Init(public, private)
	r.Mini.Init(public, private)
}
