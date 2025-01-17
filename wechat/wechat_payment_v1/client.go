package wechat_payment_v1

import (
	"errors"
)

type client struct{}

var clientVar = client{}

func Client() *client {
	return &clientVar
}

type WechatPayClient struct {
	AppId     string
	MchId     string
	SecretKey string
}

var wpc *WechatPayClient

// New  返回一个新的支付Client
func (rec *client) New(params *WechatPayClient) (*WechatPayClient, error) {
	return &WechatPayClient{
		AppId:     params.AppId,
		MchId:     params.MchId,
		SecretKey: params.SecretKey,
	}, nil
}

// Init  初始化微信支付（单商户支付方式可在应用启动时全局初始化一次）
func (rec *client) Init(params *WechatPayClient) (err error) {
	if wpc != nil {
		err = errors.New("微信支付Client已经初始化过，请勿重复初始化，如需新的支付Client请使用New方法")
		return
	}
	wpc = &WechatPayClient{
		AppId:     params.AppId,
		MchId:     params.MchId,
		SecretKey: params.SecretKey,
	}
	return
}

// Get 获取已初始化的client
func (rec *client) Get() (*WechatPayClient, error) {
	if wpc == nil {
		return nil, errors.New("微信支付Client未初始化")
	}
	return wpc, nil
}

// Which 判断使用自定义client还是全局初始化client
func (rec *client) Which(newWpc ...*WechatPayClient) (wc *WechatPayClient, err error) {
	if len(newWpc) == 1 && newWpc[0] != nil {
		return newWpc[0], nil
	}
	return rec.Get()
}
