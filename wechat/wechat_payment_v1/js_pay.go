package wechat_payment_v1

import (
	"encoding/json"
	"strconv"
	"time"
)

type jsPay struct{}

var jsPayVar = jsPay{}

func JsPay() *jsPay {
	return &jsPayVar
}

type JsPayParams struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// GetJsPay 获取JSAPI支付参数
func (rec *jsPay) GetJsPay(params *UnifiedOrderResp, newWpc ...*WechatPayClient) (js *JsPayParams, err error) {
	wc, err := Client().Which(newWpc)
	if err != nil {
		return
	}
	js = &JsPayParams{
		AppId:     params.AppId,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  params.NonceStr,
		Package:   "prepay_id=" + params.PrepayId,
		SignType:  "MD5",
	}
	paramsJson, _ := json.Marshal(js)
	js.PaySign = Utils().GetSign(paramsJson, wc.SecretKey)
	return
}
