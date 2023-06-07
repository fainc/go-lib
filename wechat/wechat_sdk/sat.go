package wechat_sdk

import (
	"sync"
)

type sat struct {
	sdk *SdkClient
}

func Sat(sdk *SdkClient) *sat {
	return &sat{sdk}
}

// Get  获取服务端接口调用凭据
func (rec *sat) Get() (token string, err error) {
	token, err = rec.getSatToken()
	if err != nil {
		return
	}
	if token == "" {
		token, err = rec.Refresh()
		return
	}
	return
}
func (rec *sat) getSatToken() (token string, err error) {
	if rec.sdk.SatRwLock == nil {
		rec.sdk.SatRwLock = new(sync.RWMutex)
	}
	defer rec.sdk.SatRwLock.RUnlock()
	rec.sdk.SatRwLock.RLock()
	token, err = Cache().GetCache("sat", rec.sdk.AppId)
	return
}

// Refresh  刷新接口调用凭据
func (rec *sat) Refresh() (token string, err error) {
	if rec.sdk.SatRwLock == nil {
		rec.sdk.SatRwLock = new(sync.RWMutex)
	}
	rec.sdk.SatRwLock.Lock()
	defer rec.sdk.SatRwLock.Unlock()
	s, err := Api().GetSat(rec.sdk)
	if err != nil {
		return
	}
	token = s.AccessToken
	err = Cache().SetCache("sat", rec.sdk.AppId, token, s.ExpiresIn)
	return
}
