package wechat_sdk

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
)

type jsSdk struct{}

var jsSdkVar = jsSdk{}

func JsSdk() *jsSdk {
	return &jsSdkVar
}

func (rec *jsSdk) GetJsApiTicket(sdk *SdkClient) (jat string, err error) {
	switch Cache().GetEngine() {
	case "redis":
		jat, err = Cache().GetRedisCache("jat", sdk.AppId)
		if err != nil {
			return
		}
	default:
		jat, err = Cache().GetMemoryCache("jat", sdk.AppId)
		if err != nil {
			return
		}
	}
	if jat == "" {
		jat, err = rec.RefreshJsApiTicket(sdk)
		return
	}
	return
}

func (rec *jsSdk) RefreshJsApiTicket(sdk *SdkClient) (jat string, err error) {
	s, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	ticket, err := Api().GetJsApiTicket(s)
	if err != nil {
		return
	}
	jat = ticket.Ticket
	switch Cache().GetEngine() {
	case "redis":
		err = Cache().SetRedisCache("jat", sdk.AppId, jat, ticket.ExpiresIn)
		if err != nil {
			return
		}
	default:
		err = Cache().SetMemoryCache("jat", sdk.AppId, jat, ticket.ExpiresIn)
		if err != nil {
			return
		}
	}
	return
}

type JsApiConfigResp struct {
	AppId     string `json:"appId"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func (rec *jsSdk) GetJsApiConfig(sdk *SdkClient, url string) (res *JsApiConfigResp, err error) {
	ticket, err := JsSdk().GetJsApiTicket(sdk)
	if err != nil {
		return
	}
	res = &JsApiConfigResp{
		AppId:     sdk.AppId,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  Utils().GetNonceStr(),
	}
	str := "jsapi_ticket=" + ticket + "&noncestr=" + res.NonceStr + "&timestamp=" + res.Timestamp + "&url=" + url
	h := sha1.New()
	h.Write([]byte(str))
	res.Signature = hex.EncodeToString(h.Sum(nil))
	return
}
