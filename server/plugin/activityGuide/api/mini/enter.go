package mini

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/service"

var (
	Guide       = new(guideApi)
	svcGuide    = service.Service.Guide
)
