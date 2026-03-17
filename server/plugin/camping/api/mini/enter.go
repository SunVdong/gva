package mini

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/service"

var (
	Site        = new(siteApi)
	Slot        = new(slotApi)
	Reservation = new(reservationApi)

	svcVenue             = service.Service.Venue
	svcVenueOpenTime     = service.Service.VenueOpenTime
	svcVenueTimeslot     = service.Service.VenueTimeslot
	svcVenueCalendar     = service.Service.VenueCalendar
	svcReservation       = service.Service.Reservation
	svcReservationReview = service.Service.ReservationReview
)
