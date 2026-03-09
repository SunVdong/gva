package request

// VenueOpenTimeBody 场地开放时间（创建/更新）
type VenueOpenTimeBody struct {
	VenueID   uint   `json:"venueId" form:"venueId" binding:"required"`
	WeekDay   int    `json:"weekDay" form:"weekDay" binding:"required,min=1,max=7"`
	OpenTime  string `json:"openTime" form:"openTime" binding:"required"`  // 如 09:00
	CloseTime string `json:"closeTime" form:"closeTime" binding:"required"` // 如 18:00
}
