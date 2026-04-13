package feedback

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FeedbackRouter struct{}

func (r *FeedbackRouter) InitFeedbackRouter(Router *gin.RouterGroup) {
	g := Router.Group("feedback").Use(middleware.OperationRecord())
	gWithout := Router.Group("feedback")
	{
		g.DELETE("deleteFeedback", feedbackApi.DeleteFeedback)
		g.DELETE("deleteFeedbackByIds", feedbackApi.DeleteFeedbackByIds)
	}
	{
		gWithout.GET("findFeedback", feedbackApi.FindFeedback)
		gWithout.GET("getFeedbackList", feedbackApi.GetFeedbackList)
	}
}
