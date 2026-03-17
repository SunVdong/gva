package service

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	ticketRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
)

type orderReview struct{}

// CreateReview 对已核销的订单发表评价（一单一评，仅本人）
func (s *orderReview) CreateReview(req ticketRequest.CreateOrderReviewRequest, userID uint) (model.OrderReview, error) {
	var order model.TicketOrder
	if err := global.GVA_DB.Where("id = ?", req.OrderID).First(&order).Error; err != nil || order.ID == 0 {
		return model.OrderReview{}, fmt.Errorf("订单不存在")
	}
	if order.UserID != userID {
		return model.OrderReview{}, fmt.Errorf("无权对该订单评价")
	}
	if order.Status != 2 {
		return model.OrderReview{}, fmt.Errorf("仅已核销订单可评价")
	}
	if order.VerifiedAt == nil {
		return model.OrderReview{}, fmt.Errorf("仅核销后的订单可评价")
	}
	var exist model.OrderReview
	if err := global.GVA_DB.Where("order_id = ?", req.OrderID).First(&exist).Error; err == nil && exist.ID != 0 {
		return model.OrderReview{}, fmt.Errorf("该订单已评价过，可先删除再重新评价")
	}
	review := model.OrderReview{
		OrderID: req.OrderID,
		UserID:  userID,
		Rating:  req.Rating,
		Content: req.Content,
	}
	return review, global.GVA_DB.Create(&review).Error
}

// DeleteReview 删除评价（仅本人可删）
func (s *orderReview) DeleteReview(reviewID uint, userID uint) error {
	var review model.OrderReview
	if err := global.GVA_DB.Where("id = ?", reviewID).First(&review).Error; err != nil || review.ID == 0 {
		return fmt.Errorf("评价不存在")
	}
	if review.UserID != userID {
		return fmt.Errorf("无权删除该评价")
	}
	return global.GVA_DB.Delete(&review).Error
}

// GetByOrderID 根据订单ID获取评价（用于详情页展示）
func (s *orderReview) GetByOrderID(orderID uint) (model.OrderReview, error) {
	var review model.OrderReview
	err := global.GVA_DB.Where("order_id = ?", orderID).First(&review).Error
	return review, err
}
