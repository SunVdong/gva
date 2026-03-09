package router

import "github.com/gin-gonic/gin"

var Router = new(router)

type router struct {
	Site            siteRouter
	TimeSlot        timeSlotRouter
	Reservation     reservationRouter
	VenueOpenTime   venueOpenTimeRouter
	VenueCalendar   venueCalendarRouter
}

// Init 初始化露营插件路由
func (r *router) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	r.Site.Init(public, private)
	r.TimeSlot.Init(public, private)
	r.Reservation.Init(public, private)
	r.VenueOpenTime.Init(public, private)
	r.VenueCalendar.Init(public, private)
}
