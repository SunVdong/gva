package service

var Service = new(service)

type service struct {
	Site        site
	TimeSlot    timeSlot
	Reservation reservation
}
