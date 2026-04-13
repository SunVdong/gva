package feedback

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	FeedbackRouter
}

var feedbackApi = api.ApiGroupApp.FeedbackApiGroup.FeedbackApi
