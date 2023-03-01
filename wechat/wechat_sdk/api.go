package wechat_sdk

import (
	"errors"
)

type api struct{}

var apiVar = api{}

func Api() *api {
	return &apiVar
}

type ApiSatResp struct {
	WxCommonResp
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (rec *api) GetSat(sdk *SdkClient) (res *ApiSatResp, err error) {
	res = &ApiSatResp{}
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + sdk.AppId + "&secret=" + sdk.Secret
	err = Request().Get(url, res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.AccessToken == "" {
		err = errors.New("获取Sat失败")
		return
	}
	return
}

type MpUserAccessTokenResp struct {
	WxCommonResp
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	OpenId         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotUser int    `json:"is_snapshotuser"`
	UnionId        string `json:"unionid"`
}

func (rec *api) GetMpUserAccessToken(sdk *SdkClient, code string) (res *MpUserAccessTokenResp, err error) {
	res = &MpUserAccessTokenResp{}
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + sdk.AppId + "&secret=" + sdk.Secret + "&code=" + code + "&grant_type=authorization_code"
	err = Request().Get(url, res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.AccessToken == "" || res.OpenId == "" {
		err = errors.New("获取用户AccessToken失败")
		return
	}
	return
}

type MpUserInfoResp struct {
	WxCommonResp
	OpenId     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	UnionId    string `json:"unionid"`
}

func (rec *api) GetMpUserInfo(accessToken string, openid string, lang string) (res *MpUserInfoResp, err error) {
	if lang == "" {
		lang = "zh_CN"
	}
	res = &MpUserInfoResp{}
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openid + "&lang=" + lang
	err = Request().Get(url, res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.OpenId == "" {
		err = errors.New("获取用户信息失败")
		return
	}
	return
}

type JsApiTicketResp struct {
	WxCommonResp
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

func (rec *api) GetJsApiTicket(accessToken string) (res *JsApiTicketResp, err error) {

	res = &JsApiTicketResp{}
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi"
	err = Request().Get(url, res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.Ticket == "" {
		err = errors.New("获取JsApiTicket失败")
		return
	}
	return
}
