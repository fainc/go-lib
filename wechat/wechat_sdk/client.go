package wechat_sdk

import (
	"errors"
	"sync"
)

type client struct{}

var clientVar = client{}

func Client() *client {
	return &clientVar
}

var clients = make(map[string]*SdkClient)

type SdkClient struct {
	AppId     string        `json:"appId"` // 小程序/公众号app_id/企业微信cropId
	Secret    string        `json:"secret"`
	SatRwLock *sync.RWMutex // Sat读写锁
	JatRwLock *sync.RWMutex // Jat读写锁
}

// Set 内存维护账号密码，如密钥变更重设密钥
func (c *client) Set(params SdkClient) (*SdkClient, error) {
	if clients[params.AppId] == nil {
		clients[params.AppId] = &SdkClient{
			AppId:     params.AppId,
			Secret:    params.Secret,
			SatRwLock: new(sync.RWMutex),
			JatRwLock: new(sync.RWMutex),
		}
	} else {
		if clients[params.AppId].Secret != params.Secret {
			clients[params.AppId].Secret = params.Secret
		}
	}
	return clients[params.AppId], nil
}
func (c *client) Get(AppId string) (sdk *SdkClient, err error) {
	if clients[AppId] != nil {
		return clients[AppId], nil
	}
	return nil, errors.New("未找到APPID " + AppId + " 的配置信息")
}
