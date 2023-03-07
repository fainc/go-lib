package wechat_sdk

type miniProgramMessage struct {
	Sdk *SdkClient
}

func MiniProgramMessage(AppId string) (*miniProgramMessage, error) {
	sdk, err := Client().Get(AppId)
	if err != nil {
		return nil, err
	}
	return &miniProgramMessage{Sdk: sdk}, nil
}

// SendUniformMessage 小程序发送统一服务消息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/uniform-message/sendUniformMessage.html
func (rec *miniProgramMessage) SendUniformMessage(p *UniformMessageParams) (res *WxCommonRes, err error) {
	token, err := Sat().Get(rec.Sdk)
	if err != nil {
		return
	}
	res, err = Api().SendUniformMessage(token, p)
	return
}

// SendSubscribeMessage 发送小程序订阅消息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/sendMessage.html
func (rec *miniProgramMessage) SendSubscribeMessage(p *MiniProgramSubscribeMessageParams) (res *WxCommonRes, err error) {
	token, err := Sat().Get(rec.Sdk)
	if err != nil {
		return
	}
	res, err = Api().SendSubscribeMessage(token, p)
	return
}
