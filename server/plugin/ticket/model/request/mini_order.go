package request

// MiniOrderCreate 小程序-创建订单（userId 由 x-token 解析注入，不从前端接收）
type MiniOrderCreate struct {
	BookerName  string `json:"bookerName" binding:"required"`
	BookerPhone string `json:"bookerPhone" binding:"required"`
	SkuID       uint   `json:"skuId" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	VisitDate   string `json:"visitDate" binding:"required"` // YYYY-MM-DD
}
