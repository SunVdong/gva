package mini

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	payConfig "github.com/silenceper/wechat/v2/pay/config"
	"github.com/silenceper/wechat/v2/pay/notify"
	"github.com/silenceper/wechat/v2/pay/order"
	"github.com/silenceper/wechat/v2/util"
)

// JSAPIParams 小程序 wx.requestPayment 所需参数（与微信文档一致）
type JSAPIParams struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// MockPaySign 模拟支付时返回的 paySign，前端可据此判断为模拟、不调 wx.requestPayment
const MockPaySign = "MOCK_SIMULATION"

// CreateJSAPI 创建 JSAPI 预支付，返回小程序调起支付所需参数。
// 当 miniprogram.mch-id / pay-key / notify-url 未配置时返回模拟参数（PaySign 为 MOCK_SIMULATION），便于本地联调。
// outTradeNo: 商户订单号（唯一）；totalFeeFen: 金额（分）；body: 商品描述；openID: 用户 openid；clientIP: 用户 IP。
func CreateJSAPI(outTradeNo string, totalFeeFen int64, body, openID, clientIP string) (*JSAPIParams, error) {
	cfg := getPayConfig()
	if cfg == nil {
		// 未配置支付时返回模拟参数，前端可根据 PaySign == MockPaySign 提示“模拟支付”并跳过 wx.requestPayment
		return &JSAPIParams{
			TimeStamp: "0",
			NonceStr:  "mock",
			Package:   "prepay_id=mock",
			SignType:  "MD5",
			PaySign:   MockPaySign,
		}, nil
	}
	if totalFeeFen <= 0 {
		return nil, fmt.Errorf("支付金额必须大于 0")
	}
	payClient := order.NewOrder(cfg)
	p := &order.Params{
		TotalFee:   strconv.FormatInt(totalFeeFen, 10),
		CreateIP:   clientIP,
		Body:       body,
		OutTradeNo: outTradeNo,
		OpenID:     openID,
		TradeType:  "JSAPI",
		SignType:   util.SignTypeMD5,
		NotifyURL:  cfg.NotifyURL,
	}
	cfgRet, err := payClient.BridgeConfig(p)
	if err != nil {
		return nil, fmt.Errorf("调起微信支付失败: %w", err)
	}
	return &JSAPIParams{
		TimeStamp: cfgRet.Timestamp,
		NonceStr:  cfgRet.NonceStr,
		Package:   cfgRet.Package,
		SignType:  cfgRet.SignType,
		PaySign:   cfgRet.PaySign,
	}, nil
}

// PaidNotifyResult 支付成功回调解析结果，供业务层根据 OutTradeNo 更新订单
type PaidNotifyResult struct {
	OutTradeNo    string
	TransactionID string
	TotalFee      int
	Attach        string
}

// ParseAndVerifyPaidNotify 解析并验签微信支付成功回调 body，返回订单号等信息；验签失败返回 error。
func ParseAndVerifyPaidNotify(body []byte) (*PaidNotifyResult, error) {
	cfg := getPayConfig()
	if cfg == nil {
		return nil, fmt.Errorf("微信支付未配置")
	}
	n := notify.NewNotify(cfg)
	var res notify.PaidResult
	if err := xml.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析回调 XML 失败: %w", err)
	}
	if res.ReturnCode == nil || *res.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("return_code 非 SUCCESS")
	}
	if res.ResultCode != nil && *res.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("result_code 非 SUCCESS")
	}
	if !n.PaidVerifySign(res) {
		return nil, fmt.Errorf("签名验证失败")
	}
	outTradeNo := ""
	if res.OutTradeNo != nil {
		outTradeNo = *res.OutTradeNo
	}
	transactionID := ""
	if res.TransactionID != nil {
		transactionID = *res.TransactionID
	}
	totalFee := 0
	if res.TotalFee != nil {
		totalFee = *res.TotalFee
	}
	attach := ""
	if res.Attach != nil {
		attach = *res.Attach
	}
	return &PaidNotifyResult{
		OutTradeNo:    outTradeNo,
		TransactionID: transactionID,
		TotalFee:      totalFee,
		Attach:        attach,
	}, nil
}

// paidNotifyResp 微信要求的回调响应根节点为 xml
type paidNotifyResp struct {
	XMLName    struct{} `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// RespondPaidNotifySuccess 返回给微信的成功 XML 响应体（写回 ResponseWriter 后不要再写 body）
func RespondPaidNotifySuccess(w io.Writer) error {
	resp := paidNotifyResp{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}
	return xml.NewEncoder(w).Encode(resp)
}

func getPayConfig() *payConfig.Config {
	c := &global.GVA_CONFIG.Miniprogram
	if c.MchID == "" || c.PayKey == "" || c.NotifyURL == "" {
		return nil
	}
	return &payConfig.Config{
		AppID:     c.AppID,
		MchID:     c.MchID,
		Key:       c.PayKey,
		NotifyURL: c.NotifyURL,
	}
}
