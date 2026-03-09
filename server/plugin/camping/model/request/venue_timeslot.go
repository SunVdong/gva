package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// VenueTimeslotSearch 场地时间段搜索
type VenueTimeslotSearch struct {
	VenueID *uint `json:"venueId" form:"venueId"`
	request.PageInfo
}
