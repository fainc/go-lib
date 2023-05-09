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
var clientRwLock = new(sync.RWMutex)

type SdkClient struct {
	AppId     string        `json:"appId"` // 小程序/公众号app_id/企业微信cropId
	Secret    string        `json:"secret"`
	SatRwLock *sync.RWMutex // Sat读写锁
	JatRwLock *sync.RWMutex // Jat读写锁
}

// New 内存维护账号密码
func (c *client) New(params SdkClient) (result *SdkClient, err error) {
	result = c.read(params.AppId)
	if result != nil {
		if result.Secret != params.Secret {
			c.update(params)
		}
		result = c.read(params.AppId)
		return
	}
	c.set(params)
	result = c.read(params.AppId)
	return
}
func (c *client) read(appId string) (result *SdkClient) {
	defer clientRwLock.RUnlock()
	clientRwLock.RLock()
	result, ok := clients[appId]
	if !ok {
		return nil
	}
	return result
}
func (c *client) set(params SdkClient) {
	defer clientRwLock.Unlock()
	clientRwLock.Lock()
	clients[params.AppId] = &SdkClient{
		AppId:     params.AppId,
		Secret:    params.Secret,
		SatRwLock: new(sync.RWMutex),
		JatRwLock: new(sync.RWMutex),
	}
}
func (c *client) update(params SdkClient) {
	clientRwLock.Lock()
	clients[params.AppId].Secret = params.Secret
	clientRwLock.Unlock()
}
func (c *client) Get(AppId string) (sdk *SdkClient, err error) {
	if clients[AppId] != nil {
		return clients[AppId], nil
	}
	return nil, errors.New("未找到APPID " + AppId + " 的配置信息")
}
