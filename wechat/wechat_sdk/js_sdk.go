package wechat_sdk

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type jsSdk struct{}

var jsSdkVar = jsSdk{}

func JsSdk() *jsSdk {
	return &jsSdkVar
}

func (rec *jsSdk) GetJsApiTicket(sdk *SdkClient) (ticket string, err error) {
	ticket, err = rec.getJsApiTicket(sdk)
	if err != nil {
		return
	}
	if ticket == "" {
		ticket, err = rec.RefreshJsApiTicket(sdk)
		return
	}
	return
}
func (rec *jsSdk) getJsApiTicket(sdk *SdkClient) (jat string, err error) {
	if sdk.JatRwLock == nil {
		sdk.JatRwLock = new(sync.RWMutex)
	}
	defer sdk.JatRwLock.RUnlock()
	sdk.JatRwLock.RLock()
	switch Cache().GetEngine() {
	case "redis":
		jat, err = Cache().GetRedisCache("jat", sdk.AppId)
		if err != nil {
			return
		}
	case "remote": // 通过远程凭据中心获取
		jat, err = RemoteCredentials().GetJat(sdk)
		if err != nil {
			return
		}
	default:
		jat, err = Cache().GetMemoryCache("jat", sdk.AppId)
		if err != nil {
			return
		}
	}
	return
}
func (rec *jsSdk) RefreshJsApiTicket(sdk *SdkClient) (jat string, err error) {
	if sdk.JatRwLock == nil {
		sdk.JatRwLock = new(sync.RWMutex)
	}
	sdk.JatRwLock.Lock()
	defer sdk.JatRwLock.Unlock()
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
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, res.NonceStr, res.Timestamp, url)
	h := sha1.New()
	h.Write([]byte(str))
	res.Signature = hex.EncodeToString(h.Sum(nil))
	return
}
