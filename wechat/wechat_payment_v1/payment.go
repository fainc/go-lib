package wechat_payment_v1

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

type payment struct{}

var paymentVar = payment{}

func Payment() *payment {
	return &paymentVar
}

type UnifiedOrderParams struct {
	TradeType     string `json:"trade_type" xml:"trade_type"`                             // * JSAPI -JSAPI支付 NATIVE -Native支付 APP -APP支付
	NotifyUrl     string `json:"notify_url" xml:"notify_url"`                             // * 公网异步通知地址 不能携带参数
	OpenId        string `json:"openid" xml:"openid"`                                     // * 用户标识
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`                         // * 商户订单号
	TotalFee      string `json:"total_fee" xml:"total_fee"`                               // * 订单金额，按分表示
	Body          string `json:"body" xml:"body"`                                         // 可选，订单说明，不传默认为 订单_${OutTradeNo}
	LimitPay      string `json:"limit_pay,omitempty" xml:"limit_pay,omitempty"`           // 可选，上传此参数no_credit--可限制用户不能使用信用卡支付
	ProfitSharing string `json:"profit_sharing,omitempty" xml:"profit_sharing,omitempty"` // 可选，Y-是，需要分账 N-否，不分账
	ProductId     string `json:"product_id,omitempty" xml:"product_id,omitempty"`         // 可选，trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	FeeType       string `json:"fee_Type,omitempty" xml:"fee_Type,omitempty"`             // 可选，符合ISO 4217标准的三位字母代码，默认人民币：CNY
	Attach        string `json:"attach,omitempty" xml:"attach,omitempty"`                 // 可选，附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	TimeExpire    string `json:"time_expire,omitempty" xml:"time_expire,omitempty"`       // 可选，订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010
}

type unifiedOrderRequestParams struct {
	XMLName        xml.Name `xml:"xml" json:"-"`
	Appid          string   `json:"appid" xml:"appid"`
	MchId          string   `json:"mch_id" xml:"mch_id"`
	NonceStr       string   `json:"nonce_str" xml:"nonce_str"`
	SpbillCreateIp string   `json:"spbill_create_ip" xml:"spbill_create_ip"`
	Sign           string   `json:"sign" xml:"sign"`
	*UnifiedOrderParams
}
type UnifiedOrderResp struct {
	AppId      string `json:"appid" xml:"appid"`
	MchId      string `json:"mch_id" xml:"mch_id"`
	NonceStr   string `json:"nonce_str" xml:"nonce_str"`
	OpenId     string `json:"openid" xml:"openid"`
	Sign       string `json:"sign" xml:"sign"`
	ResultCode string `json:"result_code" xml:"result_code"`
	ErrorCode  string `json:"err_code" xml:"err_code"`
	ErrCodeDes string `json:"err_code_des" xml:"err_code_des"`
	PrepayId   string `json:"prepay_id" xml:"prepay_id"`
	TradeType  string `json:"trade_type" xml:"trade_type"`
	CodeUrl    string `json:"code_url,omitempty" xml:"code_url,omitempty"`
	DeviceInfo string `json:"device_info,omitempty" xml:"device_info,omitempty"`
}

// UnifiedOrder 支付统一下单接口
func (rec *payment) UnifiedOrder(params *UnifiedOrderParams, newWpc ...*WechatPayClient) (resp *UnifiedOrderResp, err error) {
	wc, err := Client().Which(newWpc...)
	if err != nil {
		return
	}
	p := &unifiedOrderRequestParams{
		Appid:              wc.AppId,
		MchId:              wc.MchId,
		NonceStr:           Utils().GetNonceStr(),
		SpbillCreateIp:     "127.0.0.1",
		UnifiedOrderParams: params,
	}
	if p.Body == "" {
		p.Body = "订单_" + p.OutTradeNo
	}
	paramsJson, _ := json.Marshal(p)
	p.Sign = Utils().GetSign(paramsJson, wc.SecretKey)
	respBody, err := Request().Send("https://api.mch.weixin.qq.com/pay/unifiedorder", p)
	if err != nil {
		return
	}
	resp = &UnifiedOrderResp{}
	err = xml.Unmarshal(respBody, resp)
	if err != nil {
		return
	}
	if resp.ResultCode != "SUCCESS" {
		err = errors.New(resp.ErrCodeDes)
		return
	}
	return
}

type RefundRequestParams struct {
	TransactionId string `json:"transaction_id,omitempty" xml:"transaction_id"`     // * 原微信支付订单号（二选一A）
	OutTradeNo    string `json:"out_trade_no,omitempty" xml:"out_trade_no"`         // * 原商户订单号（二选一B）
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no"`                 // * 商户退款单号
	TotalFee      string `json:"total_fee" xml:"total_fee"`                         // * 订单金额，按分表示
	RefundFee     string `json:"refund_fee" xml:"refund_fee"`                       // * 退款金额，按分表示
	RefundDesc    string `json:"refund_desc,omitempty" xml:"refund_desc,omitempty"` // 退款原因
	NotifyUrl     string `json:"notify_url,omitempty" xml:"notify_url,omitempty"`   // 退款结果通知url
}
type refundRequestParams struct {
	XMLName  xml.Name `xml:"xml" json:"-"`
	Appid    string   `json:"appid" xml:"appid"`
	MchId    string   `json:"mch_id" xml:"mch_id"`
	NonceStr string   `json:"nonce_str" xml:"nonce_str"`
	Sign     string   `json:"sign" xml:"sign"`
	*RefundRequestParams
}

type RefundResp struct {
	ReturnCode    string `json:"return_code" xml:"return_code"`
	ReturnMsg     string `json:"return_msg" xml:"return_msg"`
	Appid         string `json:"appid" xml:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str"`
	Sign          string `json:"sign" xml:"sign"`
	ResultCode    string `json:"result_code" xml:"result_code"`
	ErrCodeDes    string `json:"err_code_des" xml:"err_code_des"`
	TransactionId string `json:"transaction_id" xml:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no"`
	RefundId      string `json:"refund_id" xml:"refund_id"`
	RefundFee     string `json:"refund_fee" xml:"refund_fee"`
}

func (rec *payment) Refund(params *RefundRequestParams, newWpc ...*WechatPayClient) (resp *RefundResp, err error) {
	wc, err := Client().Which(newWpc...)
	if err != nil {
		return
	}
	p := &refundRequestParams{
		Appid:               wc.AppId,
		MchId:               wc.MchId,
		NonceStr:            Utils().GetNonceStr(),
		RefundRequestParams: params,
	}
	paramsJson, _ := json.Marshal(p)
	p.Sign = Utils().GetSign(paramsJson, wc.SecretKey)
	if wc.Cert == nil {
		err = errors.New("wechat payment refund need cert")
		return
	}
	respBody, err := Request().Send("https://api.mch.weixin.qq.com/secapi/pay/refund", p, wc.Cert)
	if err != nil {
		return
	}
	resp = &RefundResp{}
	err = xml.Unmarshal(respBody, resp)
	if err != nil {
		return
	}
	if resp.ResultCode != "SUCCESS" {
		err = errors.New(resp.ErrCodeDes)
		return
	}
	return
}
