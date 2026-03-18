package service

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
)

type reservationReview struct{}

// CreateReview 对已核销的预约发表评价（一单一评，仅本人）
func (s *reservationReview) CreateReview(req campingRequest.CreateReservationReviewRequest, userID uint) (model.VenueReservationReview, error) {
	var res model.VenueReservation
	if err := global.GVA_DB.Where("id = ?", req.ReservationID).First(&res).Error; err != nil || res.ID == 0 {
		return model.VenueReservationReview{}, fmt.Errorf("预约不存在")
	}
	if res.UserID != userID {
		return model.VenueReservationReview{}, fmt.Errorf("无权对该预约评价")
	}
	// 仅已核销（status=1）可评价
	if res.Status != 1 {
		return model.VenueReservationReview{}, fmt.Errorf("仅核销后的预约可评价")
	}
	if res.VerifiedAt == nil {
		return model.VenueReservationReview{}, fmt.Errorf("仅核销后的预约可评价")
	}
	var exist model.VenueReservationReview
	if err := global.GVA_DB.Where("reservation_id = ?", req.ReservationID).First(&exist).Error; err == nil && exist.ID != 0 {
		return model.VenueReservationReview{}, fmt.Errorf("该预约已评价过，可先删除再重新评价")
	}
	review := model.VenueReservationReview{
		ReservationID: req.ReservationID,
		UserID:        userID,
		Rating:        req.Rating,
		Content:       req.Content,
	}
	return review, global.GVA_DB.Create(&review).Error
}

// DeleteReview 删除评价（仅本人可删）
func (s *reservationReview) DeleteReview(reviewID uint, userID uint) error {
	var review model.VenueReservationReview
	if err := global.GVA_DB.Where("id = ?", reviewID).First(&review).Error; err != nil || review.ID == 0 {
		return fmt.Errorf("评价不存在")
	}
	if review.UserID != userID {
		return fmt.Errorf("无权删除该评价")
	}
	return global.GVA_DB.Delete(&review).Error
}

// GetByReservationID 根据预约ID获取评价（用于详情页展示）
func (s *reservationReview) GetByReservationID(reservationID uint) (model.VenueReservationReview, error) {
	var review model.VenueReservationReview
	err := global.GVA_DB.Where("reservation_id = ?", reservationID).First(&review).Error
	return review, err
}
