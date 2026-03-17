package request

// CreateOrderReviewRequest 创建订单评价请求
type CreateOrderReviewRequest struct {
	OrderID uint   `json:"orderId" form:"orderId" binding:"required"`             // 订单ID
	Rating  int    `json:"rating" form:"rating" binding:"required,min=1,max=5"`   // 评分 1-5
	Content string `json:"content" form:"content" binding:"max=50"`                // 评价内容，50字内
}
