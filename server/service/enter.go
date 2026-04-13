package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	feedbackSvc "github.com/flipped-aurora/gin-vue-admin/server/service/feedback"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	FeedbackServiceGroup feedbackSvc.ServiceGroup
}
