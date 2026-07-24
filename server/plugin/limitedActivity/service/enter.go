package service

var Service = new(service)

type service struct {
	Activity activity
	Order    activityOrder
	Banner   banner
}
