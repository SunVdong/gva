package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/service"

var (
	Api          = new(api)
	serviceGuide = service.Service.Guide
)

type api struct{ Guide guide }
