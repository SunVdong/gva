package service

var Service = new(service)

type service struct {
	Scenic         scenic
	ScenicOpenTime scenicOpenTime
	Product        ticketProduct
	Sku            ticketSku
	Audience       ticketAudience
	Rule           ticketRule
	Calendar       ticketCalendar
	User           ticketUser
	Order          ticketOrder
	OrderReview    orderReview
}
