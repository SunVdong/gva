package config

// Miniprogram 微信小程序配置（用于 code2session 登录）
type Miniprogram struct {
	AppID     string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`           // 小程序 AppID
	AppSecret string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"` // 小程序 AppSecret
}
