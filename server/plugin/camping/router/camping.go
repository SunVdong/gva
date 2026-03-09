package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/api"
	"github.com/gin-gonic/gin"
)

var (
	apiSite          = api.Api.Site
	apiSlot          = api.Api.TimeSlot
	apiResv          = api.Api.Reservation
	apiVenueOpenTime = api.Api.VenueOpenTime
	apiVenueCalendar = api.Api.VenueCalendar
)

type siteRouter struct{}
type timeSlotRouter struct{}
type reservationRouter struct{}
type venueOpenTimeRouter struct{}
type venueCalendarRouter struct{}

// InitSite 场地路由
func (r *siteRouter) Init(public, private *gin.RouterGroup) {
	// 与公告一致：先 Group 单段路径，避免 "camping/site" 整段导致匹配异常
	campingSite := private.Group("camping").Group("site")
	campingSite.Use(middleware.OperationRecord()).POST("createSite", apiSite.CreateSite)
	campingSite.Use(middleware.OperationRecord()).DELETE("deleteSite", apiSite.DeleteSite)
	campingSite.Use(middleware.OperationRecord()).DELETE("deleteSiteByIds", apiSite.DeleteSiteByIds)
	campingSite.Use(middleware.OperationRecord()).PUT("updateSite", apiSite.UpdateSite)
	campingSite.GET("findSite", apiSite.FindSite)
	campingSite.GET("getSiteList", apiSite.GetSiteList)
	public.Group("camping").Group("site").GET("getSiteListPublic", apiSite.GetSiteListPublic)
	public.Group("camping").Group("site").GET("getSiteDetailPublic", apiSite.GetSiteDetailPublic)
}

// InitTimeSlot 时段路由
func (r *timeSlotRouter) Init(public, private *gin.RouterGroup) {
	slot := private.Group("camping").Group("timeSlot")
	slot.Use(middleware.OperationRecord()).POST("createTimeSlot", apiSlot.CreateTimeSlot)
	slot.Use(middleware.OperationRecord()).DELETE("deleteTimeSlot", apiSlot.DeleteTimeSlot)
	slot.Use(middleware.OperationRecord()).DELETE("deleteTimeSlotByIds", apiSlot.DeleteTimeSlotByIds)
	slot.Use(middleware.OperationRecord()).PUT("updateTimeSlot", apiSlot.UpdateTimeSlot)
	slot.GET("findTimeSlot", apiSlot.FindTimeSlot)
	slot.GET("getTimeSlotList", apiSlot.GetTimeSlotList)
	public.Group("camping").Group("timeSlot").GET("getAllTimeSlotsPublic", apiSlot.GetAllTimeSlotsPublic)
	public.Group("camping").Group("timeSlot").GET("getTimeSlotsByVenuePublic", apiSlot.GetTimeSlotsByVenuePublic)
}

// InitReservation 预约路由
func (r *reservationRouter) Init(public, private *gin.RouterGroup) {
	resv := private.Group("camping").Group("reservation")
	public.Group("camping").Group("reservation").POST("createReservation", apiResv.CreateReservation)
	public.Group("camping").Group("reservation").GET("getReservationByVerifyCodePublic", apiResv.GetReservationByVerifyCodePublic)
	public.Group("camping").Group("reservation").GET("getReservedSlotIdsPublic", apiResv.GetReservedSlotIdsPublic)
	resv.GET("getReservation", apiResv.GetReservation)
	resv.GET("getReservationList", apiResv.GetReservationList)
	resv.Use(middleware.OperationRecord()).POST("verifyReservation", apiResv.VerifyReservation)
	resv.Use(middleware.OperationRecord()).POST("verifyReservationByCode", apiResv.VerifyReservationByCode)
	resv.Use(middleware.OperationRecord()).POST("cancelReservation", apiResv.CancelReservation)
}

// InitVenueOpenTime 场地开放时间路由
func (r *venueOpenTimeRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("camping").Group("venueOpenTime")
	g.GET("getByVenue", apiVenueOpenTime.GetByVenue)
	g.Use(middleware.OperationRecord()).POST("save", apiVenueOpenTime.Save)
}

// InitVenueCalendar 场地日历路由
func (r *venueCalendarRouter) Init(public, private *gin.RouterGroup) {
	g := private.Group("camping").Group("venueCalendar")
	g.GET("getByVenue", apiVenueCalendar.GetByVenue)
	g.Use(middleware.OperationRecord()).POST("set", apiVenueCalendar.Set)
}
