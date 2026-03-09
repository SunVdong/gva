package request

// VenueCalendarSet 设置某日是否可预约
type VenueCalendarSet struct {
	VenueID uint   `json:"venueId" form:"venueId" binding:"required"`
	Date    string `json:"date" form:"date" binding:"required"` // 2006-01-02
	Status  int    `json:"status" form:"status"`              // 1可预约 0关闭
}

// VenueCalendarQuery 查询日历
type VenueCalendarQuery struct {
	VenueID uint   `json:"venueId" form:"venueId" binding:"required"`
	Start   string `json:"start" form:"start" binding:"required"` // 2006-01-02
	End     string `json:"end" form:"end" binding:"required"`     // 2006-01-02
}
