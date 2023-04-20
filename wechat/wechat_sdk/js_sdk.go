package wechat_sdk

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/fainc/go-lib/helper/str_helper"
)

type jsSdk struct {
	sdk *SdkClient
}

func JsSdk(appId, secret string) *jsSdk {
	sdk, _ := Client().New(SdkClient{
		AppId:  appId,
		Secret: secret,
	})
	return &jsSdk{sdk: sdk}
}

func (rec *jsSdk) GetJsApiTicket() (ticket string, err error) {
	ticket, err = rec.getJsApiTicket()
	if err != nil {
		return
	}
	if ticket == "" {
		ticket, err = rec.RefreshJsApiTicket()
		return
	}
	return
}
func (rec *jsSdk) getJsApiTicket() (jat string, err error) {
	if rec.sdk.JatRwLock == nil {
		rec.sdk.JatRwLock = new(sync.RWMutex)
	}
	defer rec.sdk.JatRwLock.RUnlock()
	rec.sdk.JatRwLock.RLock()
	switch Cache().GetEngine() {
	case "redis":
		jat, err = Cache().GetRedisCache("jat", rec.sdk.AppId)
		if err != nil {
			return
		}
	default:
		jat, err = Cache().GetMemoryCache("jat", rec.sdk.AppId)
		if err != nil {
			return
		}
	}
	return
}
func (rec *jsSdk) RefreshJsApiTicket() (jat string, err error) {
	if rec.sdk.JatRwLock == nil {
		rec.sdk.JatRwLock = new(sync.RWMutex)
	}
	rec.sdk.JatRwLock.Lock()
	defer rec.sdk.JatRwLock.Unlock()
	s, err := Sat(rec.sdk).Get()
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
		err = Cache().SetRedisCache("jat", rec.sdk.AppId, jat, ticket.ExpiresIn)
		if err != nil {
			return
		}
	default:
		err = Cache().SetMemoryCache("jat", rec.sdk.AppId, jat, ticket.ExpiresIn)
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

func (rec *jsSdk) GetJsApiConfig(url string) (res *JsApiConfigResp, err error) {
	ticket, err := JsSdk(rec.sdk.AppId, rec.sdk.AppId).GetJsApiTicket()
	if err != nil {
		return
	}
	res = &JsApiConfigResp{
		AppId:     rec.sdk.AppId,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  str_helper.NonceStr(),
	}
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, res.NonceStr, res.Timestamp, url)
	h := sha1.New()
	h.Write([]byte(str))
	res.Signature = hex.EncodeToString(h.Sum(nil))
	return
}
