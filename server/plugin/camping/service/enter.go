package service

var Service = new(service)

type service struct {
	Venue         venue
	VenueOpenTime venueOpenTime
	VenueTimeslot venueTimeslot
	VenueCalendar venueCalendar
	Reservation   reservation
}
