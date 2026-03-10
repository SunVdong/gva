package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/service"
)

var (
	svcScenic   = service.Service.Scenic
	svcOpenTime = service.Service.ScenicOpenTime
	svcProduct  = service.Service.Product
	svcSku      = service.Service.Sku
	svcRule     = service.Service.Rule
	svcCalendar = service.Service.Calendar
	svcOrder    = service.Service.Order
)
