package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/service"

var (
	Api                  = new(api)
	serviceVenue         = service.Service.Venue
	serviceVenueOpenTime = service.Service.VenueOpenTime
	serviceVenueTimeslot = service.Service.VenueTimeslot
	serviceVenueCalendar = service.Service.VenueCalendar
	serviceResv          = service.Service.Reservation
)

type api struct {
	Site          siteApi
	TimeSlot      timeSlotApi
	VenueOpenTime venueOpenTimeApi
	VenueCalendar venueCalendarApi
	Reservation   reservationApi
}
