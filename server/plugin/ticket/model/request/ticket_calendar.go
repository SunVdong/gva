package request

// TicketRuleItem 规则项
type TicketRuleItem struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Sort    int    `json:"sort"`
}

// TicketCalendarSetItem 日历库存设置项
type TicketCalendarSetItem struct {
	SkuID     uint   `json:"skuId"`
	VisitDate string `json:"visitDate"` // YYYY-MM-DD
	Stock     int    `json:"stock"`
	Status    int    `json:"status"` // 1可售 0关闭
}

// TicketCalendarSet 批量设置日历库存
type TicketCalendarSet struct {
	List []TicketCalendarSetItem `json:"list"`
}

// TicketCalendarSearch 日历查询
type TicketCalendarSearch struct {
	SkuID      uint   `json:"skuId" form:"skuId"`
	VisitDate  string `json:"visitDate" form:"visitDate"`   // 可选，单日
	StartDate  string `json:"startDate" form:"startDate"`   // 可选，范围开始
	EndDate    string `json:"endDate" form:"endDate"`        // 可选，范围结束
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
}
