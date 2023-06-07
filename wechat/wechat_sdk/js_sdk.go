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

func JsSdk(appID, secret string) *jsSdk {
	sdk, _ := Client().New(SdkClient{
		AppId:  appID,
		Secret: secret,
	})
	return &jsSdk{sdk: sdk}
}

func (rec *jsSdk) GetJsAPITicket() (ticket string, err error) {
	ticket, err = rec.getJsAPITicket()
	if err != nil {
		return
	}
	if ticket == "" {
		ticket, err = rec.RefreshJsAPITicket()
		return
	}
	return
}
func (rec *jsSdk) getJsAPITicket() (jat string, err error) {
	if rec.sdk.JatRwLock == nil {
		rec.sdk.JatRwLock = new(sync.RWMutex)
	}
	defer rec.sdk.JatRwLock.RUnlock()
	rec.sdk.JatRwLock.RLock()
	jat, err = Cache().GetCache("jat", rec.sdk.AppId)
	return
}
func (rec *jsSdk) RefreshJsAPITicket() (jat string, err error) {
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
	err = Cache().SetCache("jat", rec.sdk.AppId, jat, ticket.ExpiresIn)
	return
}

type JsAPIConfigResp struct {
	AppID     string `json:"appId"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func (rec *jsSdk) GetJsAPIConfig(url string) (res *JsAPIConfigResp, err error) {
	ticket, err := JsSdk(rec.sdk.AppId, rec.sdk.AppId).GetJsAPITicket()
	if err != nil {
		return
	}
	res = &JsAPIConfigResp{
		AppID:     rec.sdk.AppId,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  str_helper.NonceStr(),
	}
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, res.NonceStr, res.Timestamp, url)
	h := sha1.New()
	h.Write([]byte(str))
	res.Signature = hex.EncodeToString(h.Sum(nil))
	return
}
