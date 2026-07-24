package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/service"
)

var (
	Api             = new(api)
	serviceActivity = service.Service.Activity
	serviceOrder    = service.Service.Order
	serviceBanner   = service.Service.Banner
)

type api struct {
	Activity activityApi
	Order    activityOrderApi
	Banner   bannerApi
}
