package model

// 核销场合固定枚举（存 code，UI 展示 label）
const (
	VenueZhongshanling = "zhongshanling" // 中山陵
	VenueZhaozhao      = "zhaozhao"      // 爪爪
	VenueLululand      = "lululand"      // lululand
	VenueHongshan      = "hongshan"      // 红山
)

// VenueLabels code → 展示名
var VenueLabels = map[string]string{
	VenueZhongshanling: "中山陵",
	VenueZhaozhao:      "爪爪",
	VenueLululand:      "lululand",
	VenueHongshan:      "红山",
}

// IsValidVenue 校验场合 code 是否在白名单
func IsValidVenue(code string) bool {
	_, ok := VenueLabels[code]
	return ok
}

// VenueLabel 返回场合展示名；未知 code 返回空串
func VenueLabel(code string) string {
	return VenueLabels[code]
}

// VenueOptions 返回固定场合列表（供公开接口等使用）
func VenueOptions() []map[string]string {
	order := []string{VenueZhongshanling, VenueZhaozhao, VenueLululand, VenueHongshan}
	out := make([]map[string]string, 0, len(order))
	for _, code := range order {
		out = append(out, map[string]string{
			"code":  code,
			"label": VenueLabels[code],
		})
	}
	return out
}
