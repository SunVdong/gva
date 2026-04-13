package feedback

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	FeedbackApi
}

var feedbackService = service.ServiceGroupApp.FeedbackServiceGroup.FeedbackService
