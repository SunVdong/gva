package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/service"
)

var (
	Api                  = new(api)
	serviceScenic         = service.Service.Scenic
	serviceScenicOpenTime = service.Service.ScenicOpenTime
	serviceProduct       = service.Service.Product
	serviceSku           = service.Service.Sku
	serviceRule          = service.Service.Rule
	serviceCalendar      = service.Service.Calendar
	serviceUser          = service.Service.User
	serviceOrder         = service.Service.Order
	serviceOrderReview   = service.Service.OrderReview
)

type api struct {
	Scenic         scenicApi
	ScenicOpenTime scenicOpenTimeApi
	Product        ticketProductApi
	Sku            ticketSkuApi
	Rule           ticketRuleApi
	Calendar       ticketCalendarApi
	User           ticketUserApi
	Order          ticketOrderApi
}
