package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VenueReservationReview 预约评价（仅核销后的预约可评价，一单一评）
// 评分 1-5 对应：非常差、较差、一般、推荐、超赞
type VenueReservationReview struct {
	global.GVA_MODEL
	ReservationID uint   `json:"reservationId" form:"reservationId" gorm:"column:reservation_id;comment:预约ID;uniqueIndex:idx_reservation_review;"`
	UserID        uint   `json:"userId" form:"userId" gorm:"column:user_id;comment:评价用户ID;"`
	Rating        int    `json:"rating" form:"rating" gorm:"column:rating;comment:评分1-5;"`
	Content       string `json:"content" form:"content" gorm:"column:content;comment:评价内容50字内;size:50;"`
}

// TableName 表名
func (VenueReservationReview) TableName() string {
	return "venue_reservation_reviews"
}
