package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

// VenueReservationSearch 预约搜索
type VenueReservationSearch struct {
	UserID      *uint      `json:"userId" form:"userId"`           // 小程序「我的预约」按用户筛选
	VenueID     *uint      `json:"venueId" form:"venueId"`
	ReserveDate *time.Time `json:"reserveDate" form:"reserveDate"`
	Status      *int       `json:"status" form:"status"`
	VerifyCode  string     `json:"verifyCode" form:"verifyCode"`
	request.PageInfo
}

// CreateVenueReservationRequest 创建预约请求
type CreateVenueReservationRequest struct {
	VenueID      uint   `json:"venueId" form:"venueId" binding:"required"`
	ReserveDate  string `json:"reserveDate" form:"reserveDate" binding:"required"` // 2006-01-02
	TimeslotID   uint   `json:"timeslotId" form:"timeslotId" binding:"required"`
	ContactName  string `json:"contactName" form:"contactName" binding:"required"`
	ContactPhone string `json:"contactPhone" form:"contactPhone" binding:"required"`
	ContactCount int    `json:"contactCount" form:"contactCount" binding:"required,min=1"`
	Remark       string `json:"remark" form:"remark"`
}

// UpdateVenueReservationRequest 修改预约请求
// 仅允许修改预约日期、时段和联系人信息，不支持跨场地修改
type UpdateVenueReservationRequest struct {
	ID           uint   `json:"id" form:"id" binding:"required"`                         // 预约ID
	ReserveDate  string `json:"reserveDate" form:"reserveDate" binding:"required"`       // 2006-01-02
	TimeslotID   uint   `json:"timeslotId" form:"timeslotId" binding:"required"`        // 新的时间段ID
	ContactName  string `json:"contactName" form:"contactName" binding:"required"`      // 联系人
	ContactPhone string `json:"contactPhone" form:"contactPhone" binding:"required"`    // 手机号
	ContactCount int    `json:"contactCount" form:"contactCount" binding:"required,min=1"` // 预约人数
}
