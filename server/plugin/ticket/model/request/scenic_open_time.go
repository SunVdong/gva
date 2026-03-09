package request

// ScenicOpenTimeItem 开放时间项
type ScenicOpenTimeItem struct {
	ScenicID  uint   `json:"scenicId"`
	WeekDay   int    `json:"weekDay"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
}

// ScenicOpenTimeSave 保存景区开放时间
type ScenicOpenTimeSave struct {
	ScenicID uint                 `json:"scenicId" binding:"required"`
	List     []ScenicOpenTimeItem `json:"list"`
}
