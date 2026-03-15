package config

// Miniprogram 微信小程序配置（登录 + 支付）
type Miniprogram struct {
	AppID     string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`             // 小程序 AppID
	AppSecret string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"` // 小程序 AppSecret
	// 微信支付（JSAPI，与小程序同主体）
	MchID     string `mapstructure:"mch-id" json:"mch-id" yaml:"mch-id"`             // 商户号
	PayKey    string `mapstructure:"pay-key" json:"pay-key" yaml:"pay-key"`         // 支付 API 密钥（v2 密钥）
	NotifyURL string `mapstructure:"notify-url" json:"notify-url" yaml:"notify-url"` // 支付结果回调地址，如 https://yourdomain.com/api/mini/pay/notify
}
