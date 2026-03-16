package request

// MiniOrderItem 小程序下单-订单项
type MiniOrderItem struct {
	SkuID     uint   `json:"skuId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	VisitDate string `json:"visitDate" binding:"required"` // YYYY-MM-DD
}

// MiniOrderCreate 小程序-创建订单
type MiniOrderCreate struct {
	UserID      uint            `json:"userId" binding:"required"`
	BookerName  string          `json:"bookerName" binding:"required"`
	BookerPhone string          `json:"bookerPhone" binding:"required"`
	Items       []MiniOrderItem `json:"items" binding:"required,min=1"`
}
