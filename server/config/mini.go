package config

// Miniprogram 微信小程序配置（登录 + 支付）
type Miniprogram struct {
	AppID     string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`             // 小程序 AppID
	AppSecret string `mapstructure:"app-secret" json:"app-secret" yaml:"app-secret"` // 小程序 AppSecret
	// 微信支付 APIv3（小程序 JSAPI，与小程序同主体）
	MchID              string `mapstructure:"mch-id" json:"mch-id" yaml:"mch-id"`                                           // 商户号
	APIv3Key           string `mapstructure:"api-v3-key" json:"api-v3-key" yaml:"api-v3-key"`                               // APIv3 密钥（32 位）
	MchAPIv3SerialNo   string `mapstructure:"mch-api-serial-no" json:"mch-api-serial-no" yaml:"mch-api-serial-no"`             // 商户 API 证书序列号
	MchPrivateKey      string `mapstructure:"mch-private-key" json:"mch-private-key" yaml:"mch-private-key"`                 // 商户 API 私钥 PEM（apiclient_key.pem 全文，与 file 二选一）
	MchPrivateKeyFile  string `mapstructure:"mch-private-key-file" json:"mch-private-key-file" yaml:"mch-private-key-file"` // 或填写私钥文件路径
	NotifyURL          string `mapstructure:"notify-url" json:"notify-url" yaml:"notify-url"`                                 // 支付结果回调，如 https://yourdomain.com/api/mini/pay/notify
	WxPayPublicKey     string `mapstructure:"wx-pay-public-key" json:"wx-pay-public-key" yaml:"wx-pay-public-key"`           // 微信支付公钥内容（PEM）
	WxPayPublicKeyID   string `mapstructure:"wx-pay-public-key-id" json:"wx-pay-public-key-id" yaml:"wx-pay-public-key-id"` // 微信支付公钥ID（含 PUB_KEY_ID_ 前缀）
}
