package mini

// ApiGroup 小程序端 API 分组，后续可在此扩展多个 API（如 CommonApi、AuthApi、PayApi 等）
type ApiGroup struct {
	CommonApi
	AuthApi
	PayApi
}
