package mini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

var (
	wxV3Mu     sync.Mutex
	wxV3Client *wxv3.ClientV3
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

// payConfigMissingReason 若非空则说明支付配置未就绪，便于区分「写错文件」与「缺字段」。
func payConfigMissingReason() string {
	c := &global.GVA_CONFIG.Miniprogram
	var missing []string
	if strings.TrimSpace(c.AppID) == "" {
		missing = append(missing, "app-id")
	}
	if strings.TrimSpace(c.MchID) == "" {
		missing = append(missing, "mch-id")
	}
	if strings.TrimSpace(c.APIv3Key) == "" {
		missing = append(missing, "api-v3-key")
	}
	if strings.TrimSpace(c.MchAPIv3SerialNo) == "" {
		missing = append(missing, "mch-api-serial-no")
	}
	if strings.TrimSpace(c.NotifyURL) == "" {
		missing = append(missing, "notify-url")
	}
	if strings.TrimSpace(c.WxPayPublicKey) == "" {
		missing = append(missing, "wx-pay-public-key")
	}
	if strings.TrimSpace(c.WxPayPublicKeyID) == "" {
		missing = append(missing, "wx-pay-public-key-id")
	}
	if strings.TrimSpace(c.MchPrivateKey) == "" && strings.TrimSpace(c.MchPrivateKeyFile) == "" {
		missing = append(missing, "mch-private-key 或 mch-private-key-file")
	}
	if len(missing) == 0 {
		return ""
	}
	return fmt.Sprintf(
		"缺少或未生效字段: %s。请对照启动日志里「config 的路径」检查该文件中的 miniprogram 段；Gin release 且无 config.release.yaml 时会使用 config.yaml（不会自动读 config.debug.yaml）。可通过环境变量 GVA_CONFIG 或启动参数 -c 指向含支付配置的 yaml。",
		strings.Join(missing, "、"),
	)
}

// getWxV3Client 懒加载微信 V3 客户端；仅在成功拉取平台证书后缓存。
// 初始化失败不缓存错误：网络恢复或修正配置后，下一次请求会自动重试，无需重启进程。
func getWxV3Client() (*wxv3.ClientV3, error) {
	if r := payConfigMissingReason(); r != "" {
		return nil, fmt.Errorf("微信支付未配置完整：%s", r)
	}
	wxV3Mu.Lock()
	defer wxV3Mu.Unlock()
	if wxV3Client != nil {
		return wxV3Client, nil
	}
	pemStr, err := loadMchPrivateKeyPEM(&global.GVA_CONFIG.Miniprogram)
	if err != nil {
		return nil, err
	}
	m := &global.GVA_CONFIG.Miniprogram
	cli, err := wxv3.NewClientV3(m.MchID, m.MchAPIv3SerialNo, m.APIv3Key, pemStr)
	if err != nil {
		return nil, err
	}
	if err := cli.AutoVerifySignByPublicKey([]byte(m.WxPayPublicKey), m.WxPayPublicKeyID); err != nil {
		return nil, fmt.Errorf("加载微信支付公钥失败（请检查 wx-pay-public-key 与 wx-pay-public-key-id）: %w", err)
	}
	wxV3Client = cli
	return wxV3Client, nil
}

func refundNotifyURL() string {
	if u := strings.TrimSpace(global.GVA_CONFIG.Miniprogram.RefundNotifyURL); u != "" {
		return u
	}
	u := strings.TrimSpace(global.GVA_CONFIG.Miniprogram.NotifyURL)
	if u == "" {
		return ""
	}
	if strings.Contains(u, "/pay/notify") {
		return strings.Replace(u, "/pay/notify", "/pay/refund/notify", 1)
	}
	return ""
}

// CreateJSAPI 使用微信支付 V3 创建 JSAPI/小程序预支付，返回 wx.requestPayment 所需参数。
func CreateJSAPI(outTradeNo string, totalFeeFen int64, body, openID, clientIP string) (*JSAPIParams, error) {
	if r := payConfigMissingReason(); r != "" {
		return nil, fmt.Errorf("微信支付未配置完整：%s", r)
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
	if r := payConfigMissingReason(); r != "" {
		return nil, fmt.Errorf("微信支付未配置：%s", r)
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

// RefundResult 退款结果
type RefundResult struct {
	RefundID    string // 微信退款单号
	OutRefundNo string // 商户退款单号
	Status      string // SUCCESS / CLOSED / PROCESSING / ABNORMAL
}

// CreateRefund 调用微信支付 V3 申请退款，全额退款。
func CreateRefund(transactionID, outRefundNo string, totalFen, refundFen int, reason string) (*RefundResult, error) {
	client, err := getWxV3Client()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	bm := make(gopay.BodyMap)
	bm.Set("transaction_id", transactionID).
		Set("out_refund_no", outRefundNo)
	if notifyURL := refundNotifyURL(); notifyURL != "" {
		bm.Set("notify_url", notifyURL)
	}
	if reason != "" {
		bm.Set("reason", reason)
	}
	bm.SetBodyMap("amount", func(bm gopay.BodyMap) {
		bm.Set("refund", refundFen).
			Set("total", totalFen).
			Set("currency", "CNY")
	})
	wxRsp, err := client.V3Refund(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("调起微信退款失败: %w", err)
	}
	if wxRsp.Code != wxv3.Success {
		return nil, fmt.Errorf("微信退款失败: http=%d %s", wxRsp.Code, wxRsp.Error)
	}
	if wxRsp.Response == nil {
		return nil, fmt.Errorf("微信退款未返回结果")
	}
	return &RefundResult{
		RefundID:    wxRsp.Response.RefundId,
		OutRefundNo: wxRsp.Response.OutRefundNo,
		Status:      wxRsp.Response.Status,
	}, nil
}

// RefundNotifyResult 退款回调解析结果
type RefundNotifyResult struct {
	OutRefundNo   string
	RefundID      string
	TransactionID string
	RefundStatus  string
	SuccessTime   string
}

// ParseAndVerifyRefundNotify 解析并验签微信支付 V3 退款回调。
func ParseAndVerifyRefundNotify(req *http.Request) (*RefundNotifyResult, error) {
	if r := payConfigMissingReason(); r != "" {
		return nil, fmt.Errorf("微信支付未配置：%s", r)
	}
	client, err := getWxV3Client()
	if err != nil {
		return nil, err
	}
	notifyReq, err := wxv3.V3ParseNotify(req)
	if err != nil {
		return nil, fmt.Errorf("解析回调失败: %w", err)
	}
	if err := notifyReq.VerifySignByPKMap(client.WxPublicKeyMap()); err != nil {
		return nil, fmt.Errorf("签名验证失败: %w", err)
	}
	m := &global.GVA_CONFIG.Miniprogram
	refundRes, err := notifyReq.DecryptRefundCipherText(m.APIv3Key)
	if err != nil {
		return nil, fmt.Errorf("解密退款回调数据失败: %w", err)
	}
	return &RefundNotifyResult{
		OutRefundNo:   refundRes.OutRefundNo,
		RefundID:      refundRes.RefundId,
		TransactionID: refundRes.TransactionId,
		RefundStatus:  refundRes.RefundStatus,
		SuccessTime:   refundRes.SuccessTime,
	}, nil
}
