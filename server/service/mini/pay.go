package mini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/go-pay/gopay"
	wxv3 "github.com/go-pay/gopay/wechat/v3"
)

// JSAPIParams 小程序 wx.requestPayment 所需参数（与微信文档一致，V3 signType 为 RSA）
type JSAPIParams struct {
	AppId     string `json:"appId,omitempty"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// MockPaySign 模拟支付时返回的 paySign，前端可据此判断为模拟、不调 wx.requestPayment
const MockPaySign = "MOCK_SIMULATION"

var (
	wxV3Mu      sync.Mutex
	wxV3Client  *wxv3.ClientV3
	wxV3InitErr error
)

func loadMchPrivateKeyPEM(m *config.Miniprogram) (string, error) {
	if m.MchPrivateKeyFile != "" {
		b, err := os.ReadFile(m.MchPrivateKeyFile)
		if err != nil {
			return "", fmt.Errorf("读取 mch-private-key-file 失败: %w", err)
		}
		return string(b), nil
	}
	if m.MchPrivateKey != "" {
		return m.MchPrivateKey, nil
	}
	return "", fmt.Errorf("未配置 mch-private-key 或 mch-private-key-file")
}

func isPayConfigured() bool {
	c := &global.GVA_CONFIG.Miniprogram
	if c.AppID == "" || c.MchID == "" || c.NotifyURL == "" || c.APIv3Key == "" || c.MchAPIv3SerialNo == "" {
		return false
	}
	if c.MchPrivateKey == "" && c.MchPrivateKeyFile == "" {
		return false
	}
	return true
}

func getWxV3Client() (*wxv3.ClientV3, error) {
	if !isPayConfigured() {
		return nil, fmt.Errorf("微信支付未配置完整（需 app-id、mch-id、api-v3-key、mch-api-serial-no、notify-url 及商户私钥）")
	}
	wxV3Mu.Lock()
	defer wxV3Mu.Unlock()
	if wxV3Client != nil {
		return wxV3Client, nil
	}
	if wxV3InitErr != nil {
		return nil, wxV3InitErr
	}
	pemStr, err := loadMchPrivateKeyPEM(&global.GVA_CONFIG.Miniprogram)
	if err != nil {
		wxV3InitErr = err
		return nil, err
	}
	m := &global.GVA_CONFIG.Miniprogram
	cli, err := wxv3.NewClientV3(m.MchID, m.MchAPIv3SerialNo, m.APIv3Key, pemStr)
	if err != nil {
		wxV3InitErr = err
		return nil, err
	}
	if err := cli.AutoVerifySign(); err != nil {
		wxV3InitErr = fmt.Errorf("同步微信平台证书失败（请检查商户号、证书序列号、私钥与网络）: %w", err)
		return nil, wxV3InitErr
	}
	wxV3Client = cli
	return wxV3Client, nil
}

// CreateJSAPI 使用微信支付 V3 创建 JSAPI/小程序预支付，返回 wx.requestPayment 所需参数。
// 当支付未配置完整时返回模拟参数（PaySign 为 MOCK_SIMULATION），便于本地联调。
func CreateJSAPI(outTradeNo string, totalFeeFen int64, body, openID, clientIP string) (*JSAPIParams, error) {
	if !isPayConfigured() {
		return &JSAPIParams{
			TimeStamp: "0",
			NonceStr:  "mock",
			Package:   "prepay_id=mock",
			SignType:  "RSA",
			PaySign:   MockPaySign,
		}, nil
	}
	if totalFeeFen <= 0 {
		return nil, fmt.Errorf("支付金额必须大于 0")
	}
	client, err := getWxV3Client()
	if err != nil {
		return nil, err
	}
	m := &global.GVA_CONFIG.Miniprogram
	ctx := context.Background()
	bm := make(gopay.BodyMap)
	bm.Set("appid", m.AppID).
		Set("description", body).
		Set("out_trade_no", outTradeNo).
		Set("notify_url", m.NotifyURL).
		Set("time_expire", time.Now().Add(2*time.Hour).Format(time.RFC3339))
	bm.SetBodyMap("amount", func(bm gopay.BodyMap) {
		bm.Set("total", int(totalFeeFen)).
			Set("currency", "CNY")
	})
	bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
		bm.Set("openid", openID)
	})
	if clientIP != "" {
		bm.SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.Set("payer_client_ip", clientIP)
		})
	}
	wxRsp, err := client.V3TransactionJsapi(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("调起微信支付失败: %w", err)
	}
	if wxRsp.Code != wxv3.Success {
		return nil, fmt.Errorf("微信下单失败: http=%d %s", wxRsp.Code, wxRsp.Error)
	}
	if wxRsp.Response == nil || wxRsp.Response.PrepayId == "" {
		return nil, fmt.Errorf("微信下单未返回 prepay_id")
	}
	applet, err := client.PaySignOfApplet(m.AppID, wxRsp.Response.PrepayId)
	if err != nil {
		return nil, fmt.Errorf("生成支付签名失败: %w", err)
	}
	return &JSAPIParams{
		AppId:     applet.AppId,
		TimeStamp: applet.TimeStamp,
		NonceStr:  applet.NonceStr,
		Package:   applet.Package,
		SignType:  applet.SignType,
		PaySign:   applet.PaySign,
	}, nil
}

// PaidNotifyResult 支付成功回调解析结果，供业务层根据 OutTradeNo 更新订单
type PaidNotifyResult struct {
	OutTradeNo    string
	TransactionID string
	TotalFee      int
	Attach        string
}

// ParseAndVerifyPaidNotify 解析并验签微信支付 V3 回调，解密 resource 后返回订单信息；失败返回 error。
func ParseAndVerifyPaidNotify(req *http.Request) (*PaidNotifyResult, error) {
	if !isPayConfigured() {
		return nil, fmt.Errorf("微信支付未配置")
	}
	client, err := getWxV3Client()
	if err != nil {
		return nil, err
	}
	notifyReq, err := wxv3.V3ParseNotify(req)
	if err != nil {
		return nil, fmt.Errorf("解析回调失败: %w", err)
	}
	if notifyReq.EventType != "" && notifyReq.EventType != "TRANSACTION.SUCCESS" {
		return nil, fmt.Errorf("非支付成功通知: %s", notifyReq.EventType)
	}
	if err := notifyReq.VerifySignByPKMap(client.WxPublicKeyMap()); err != nil {
		return nil, fmt.Errorf("签名验证失败: %w", err)
	}
	m := &global.GVA_CONFIG.Miniprogram
	payRes, err := notifyReq.DecryptPayCipherText(m.APIv3Key)
	if err != nil {
		return nil, fmt.Errorf("解密回调数据失败: %w", err)
	}
	if payRes.TradeState != wxv3.TradeStateSuccess {
		return nil, fmt.Errorf("交易状态非成功: %s", payRes.TradeState)
	}
	total := 0
	if payRes.Amount != nil {
		total = payRes.Amount.Total
	}
	return &PaidNotifyResult{
		OutTradeNo:    payRes.OutTradeNo,
		TransactionID: payRes.TransactionId,
		TotalFee:      total,
		Attach:        payRes.Attach,
	}, nil
}

// RespondPaidNotifySuccess 返回给微信 V3 的成功 JSON（写入后不要再写 body）
func RespondPaidNotifySuccess(w io.Writer) error {
	return json.NewEncoder(w).Encode(&wxv3.V3NotifyRsp{Code: "SUCCESS", Message: "成功"})
}
