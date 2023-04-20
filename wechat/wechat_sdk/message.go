package wechat_sdk

type message struct {
	sdk *SdkClient
}

func Message(AppId, Secret string) *message {
	sdk, _ := Client().New(SdkClient{
		AppId:  AppId,
		Secret: Secret,
	})
	return &message{sdk: sdk}
}

// SendUniformMessage 小程序发送统一服务消息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/uniform-message/sendUniformMessage.html
func (rec *message) SendUniformMessage(p *UniformMessageParams) (res *WxCommonRes, err error) {
	token, err := Sat(rec.sdk).Get()
	if err != nil {
		return
	}
	res, err = Api().SendUniformMessage(token, p)
	return
}

// SendSubscribeMessage 发送小程序订阅消息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/sendMessage.html
func (rec *message) SendSubscribeMessage(p *MiniProgramSubscribeMessageParams) (res *WxCommonRes, err error) {
	token, err := Sat(rec.sdk).Get()
	if err != nil {
		return
	}
	res, err = Api().SendSubscribeMessage(token, p)
	return
}
