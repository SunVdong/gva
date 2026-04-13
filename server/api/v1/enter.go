package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	feedbackApi "github.com/flipped-aurora/gin-vue-admin/server/api/v1/feedback"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/mini"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	FeedbackApiGroup feedbackApi.ApiGroup
	MiniApiGroup     mini.ApiGroup
}
