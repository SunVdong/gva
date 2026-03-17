package request

// CreateReservationReviewRequest 创建预约评价请求
type CreateReservationReviewRequest struct {
	ReservationID uint   `json:"reservationId" form:"reservationId" binding:"required"`       // 预约ID
	Rating        int    `json:"rating" form:"rating" binding:"required,min=1,max=5"`         // 评分 1-5
	Content       string `json:"content" form:"content" binding:"max=50"`                      // 评价内容，50字内
}
