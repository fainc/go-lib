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
	switch Cache().GetEngine() {
	case "redis":
		token, err = Cache().GetRedisCache("sat", rec.sdk.AppId)
		if err != nil {
			return
		}
	default:
		token, err = Cache().GetMemoryCache("sat", rec.sdk.AppId)
		if err != nil {
			return
		}
	}
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
	switch Cache().GetEngine() {
	case "redis":
		err = Cache().SetRedisCache("sat", rec.sdk.AppId, token, s.ExpiresIn)
		if err != nil {
			return
		}
	default:
		err = Cache().SetMemoryCache("sat", rec.sdk.AppId, token, s.ExpiresIn)
		if err != nil {
			return
		}
	}
	return
}
