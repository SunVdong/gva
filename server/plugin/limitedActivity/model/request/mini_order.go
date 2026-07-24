package request

// MiniOrderCreate 小程序-创建活动订单（userId 由 x-token 解析注入）
type MiniOrderCreate struct {
	ActivityID   uint   `json:"activityId" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
	ContactName  string `json:"contactName" binding:"required"`
	ContactPhone string `json:"contactPhone" binding:"required"`
}
