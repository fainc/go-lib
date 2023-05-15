package wechat_sdk

import (
	"errors"
)

type mpApi struct{}

var mpApiVar = mpApi{}

func MpApi() *mpApi {
	return &mpApiVar
}

type MpQrCodeParams struct {
	ExpireSeconds int    `json:"expire_seconds"`
	ActionName    string `json:"action_name"`
	ActionInfo    string `json:"action_info"`
	SceneId       int    `json:"scene_id"`
	SceneStr      string `json:"scene_str"`
}
type MpQrCodeRes struct {
	WxCommonRes
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
	Path          string `json:"path"`
}

func (rec *mpApi) GetQrCode(access string, params *MpQrCodeParams) (res *MpQrCodeRes, err error) {
	url := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + access
	err = Request().Post(url, params, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.Ticket == "" || res.Url == "" {
		err = errors.New("获取二维码失败")
		return
	}
	return
}

type MpTemplateMessageParams struct {
	ToUser      string                   `json:"touser"`
	TemplateId  string                   `json:"template_id"`
	Url         string                   `json:"url,omitempty"`
	MiniProgram *miniProgram             `json:"miniprogram,omitempty"`
	ClientMsgId string                   `json:"client_msg_id,omitempty"`
	Data        map[string]*MessageValue `json:"data"`
}
type MpTemplateMessageRes struct {
	WxCommonRes
	MsgId string `json:"msgid"`
}

func (rec *mpApi) SendMpTemplateMessage(access string, params *MpTemplateMessageParams) (res *MpTemplateMessageRes, err error) {
	url := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + access
	err = Request().Post(url, params, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	return
}
