package wechat_payment_v1

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
	"time"

	"github.com/fainc/go-lib/crypto/md5_crypto"
	"github.com/fainc/go-lib/helper/signature_helper"
	"github.com/fainc/go-lib/helper/str_helper"
)

type utils struct{}

var utilsVar = utils{}

func Utils() *utils {
	return &utilsVar
}

// GetNonceStr 获取随机字符串
func (rec *utils) GetNonceStr() string {
	return str_helper.NonceStr()
}

// GetSign 微信支付统一签名方法
func (rec *utils) GetSign(p []byte, key string) string {
	var m map[string]string
	_ = json.Unmarshal(p, &m)
	str := signature_helper.SignatureStr(m, &signature_helper.SignatureStrKeyOptions{Key: "key", Value: key}, []string{"sign", "paySign"})
	return md5_crypto.Md5(str, true)
}

type PaymentRawNotify struct {
	AppId         string `json:"appid" xml:"appid"`
	Attach        string `json:"attach" xml:"attach"`
	BankType      string `json:"bank_type" xml:"bank_type"`
	FeeType       string `json:"fee_type" xml:"fee_type"`
	IsSubscribe   string `json:"is_subscribe" xml:"is_subscribe"`
	MchId         string `json:"mch_id" xml:"mch_id"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str"`
	OpenId        string `json:"openid" xml:"openid"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`
	ResultCode    string `json:"result_code" xml:"result_code"`
	ReturnCode    string `json:"return_code" xml:"return_code"`
	Sign          string `json:"sign" xml:"sign"`
	TimeEnd       string `json:"time_end" xml:"time_end"`
	TotalFee      string `json:"total_fee" xml:"total_fee"`
	CouponFee     string `json:"coupon_fee" xml:"coupon_fee"`
	CouponCount   string `json:"coupon_count" xml:"coupon_count"`
	CouponType    string `json:"coupon_type" xml:"coupon_type"`
	CouponId      string `json:"coupon_id" xml:"coupon_id"`
	TradeType     string `json:"trade_type" xml:"trade_type"`
	TransactionId string `json:"transaction_id" xml:"transaction_id"`
}

type PaymentNotify struct {
	AppId         string  `json:"appid" xml:"appid"`
	Attach        string  `json:"attach" xml:"attach"`
	IsSubscribe   bool    `json:"is_subscribe" xml:"is_subscribe"`
	MchId         string  `json:"mch_id" xml:"mch_id"`
	OpenId        string  `json:"openid" xml:"openid"`
	OutTradeNo    string  `json:"out_trade_no" xml:"out_trade_no"`
	PayTime       string  `json:"pay_time" xml:"pay_time"`
	TotalFee      float64 `json:"total_fee" xml:"total_fee"`
	TransactionId string  `json:"transaction_id" xml:"transaction_id"`
}

// ParsePaymentNotify 解析微信支付通知数据
func (rec *utils) ParsePaymentNotify(body []byte, key string) (p *PaymentNotify, raw *PaymentRawNotify, err error) {
	raw = &PaymentRawNotify{}
	err = xml.Unmarshal(body, raw)
	if err != nil {
		return
	}
	params, err := json.Marshal(raw)
	if err != nil {
		return
	}
	sign := rec.GetSign(params, key)
	if sign != raw.Sign {
		err = errors.New("支付通知签名验证失败")
		return
	}
	if raw.ReturnCode != "SUCCESS" || raw.ResultCode != "SUCCESS" {
		err = errors.New("未成功的支付通知")
		return
	}
	float, err := strconv.ParseFloat(raw.TotalFee, 64)
	if err != nil {
		return
	}
	PayTime := raw.TimeEnd
	parse, errT := time.Parse("20060102150405", raw.TimeEnd)
	if errT == nil {
		PayTime = parse.Format("2006-01-02 15:04:05")
	}
	p = &PaymentNotify{
		AppId:         raw.AppId,
		Attach:        raw.Attach,
		IsSubscribe:   raw.IsSubscribe == "Y",
		MchId:         raw.MchId,
		OpenId:        raw.OpenId,
		OutTradeNo:    raw.OutTradeNo,
		PayTime:       PayTime,
		TotalFee:      float / 100,
		TransactionId: raw.TransactionId,
	}
	return
}

type CommonResp struct {
	ReturnCode string `xml:"return_code" json:"return_code"`
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`
}

func (rec *utils) SuccessResp() *CommonResp {
	return &CommonResp{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}
}
