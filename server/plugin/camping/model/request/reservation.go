package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

// CampingReservationSearch 预约搜索
type CampingReservationSearch struct {
	SiteID      *uint      `json:"siteId" form:"siteId"`
	ReserveDate *time.Time `json:"reserveDate" form:"reserveDate"`
	Status      *int       `json:"status" form:"status"`
	VerifyCode  string     `json:"verifyCode" form:"verifyCode"`
	request.PageInfo
}

// CreateCampingReservationRequest 创建预约请求
type CreateCampingReservationRequest struct {
	SiteID      uint   `json:"siteId" form:"siteId" binding:"required"`
	ReserveDate string `json:"reserveDate" form:"reserveDate" binding:"required"` // 格式 2006-01-02
	TimeSlotID  uint   `json:"timeSlotId" form:"timeSlotId" binding:"required"`
	BookerName  string `json:"bookerName" form:"bookerName" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	PeopleCount int    `json:"peopleCount" form:"peopleCount" binding:"required,min=1"`
	Remark      string `json:"remark" form:"remark"`
}
