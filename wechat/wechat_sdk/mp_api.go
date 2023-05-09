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
	ToUser           string `json:"expire_seconds"`
	TemplateId       string `json:"action_name"`
	Page             string `json:"action_info"`
	MiniProgramState string `json:"scene_id"`
	Lang             string `json:"scene_str"`
}
type MpQrCodeRes struct {
	WxCommonRes
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
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