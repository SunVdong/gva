package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/service"

var (
	Api          = new(api)
	serviceSite  = service.Service.Site
	serviceSlot  = service.Service.TimeSlot
	serviceResv  = service.Service.Reservation
)

type api struct {
	Site        siteApi
	TimeSlot    timeSlotApi
	Reservation reservationApi
}
