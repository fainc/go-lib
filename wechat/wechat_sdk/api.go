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
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + sdk.AppId + "&secret=" + sdk.Secret
	err = Request().Get(url, &res)
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
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + sdk.AppId + "&secret=" + sdk.Secret + "&code=" + code + "&grant_type=authorization_code"
	err = Request().Get(url, &res)
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
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openid + "&lang=" + lang
	err = Request().Get(url, &res)
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
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi"
	err = Request().Get(url, &res)
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

type Code2SessionResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	WxCommonResp
}

func (rec *api) Code2Session(sdk *SdkClient, code string) (res *Code2SessionResp, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + sdk.AppId + "&secret=" + sdk.Secret + "&js_code=" + code + "&grant_type=authorization_code"
	err = Request().Get(url, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.OpenId == "" || res.SessionKey == "" {
		err = errors.New("获取用户Session失败")
		return
	}
	return
}

type PaidUnionIdResp struct {
	UnionId string `json:"unionid"`
	WxCommonResp
}

func (rec *api) GetPaidUnionId(access string, p *GetPaidUnionIdParams) (res *PaidUnionIdResp, err error) {
	url := "https://api.weixin.qq.com/wxa/getpaidunionid?access_token=" + access + "&openid=" + p.OpenId
	if p.TransactionId != "" {
		url = url + "&transaction_id=" + p.TransactionId
	}
	if p.TransactionId == "" && p.OutTradeNo != "" {
		url = url + "&out_trade_no=" + p.OutTradeNo + "&mch_id=" + p.MchId
	}
	err = Request().Get(url, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.UnionId == "" {
		err = errors.New("获取用户UnionId失败")
		return
	}
	return
}

type UserPhoneNumberResp struct {
	WxCommonResp
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     int    `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			AppId     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

func (rec *api) GetUserPhoneNumber(access string, code string) (res *UserPhoneNumberResp, err error) {
	url := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + access
	p := make(map[string]string)
	p["code"] = code
	err = Request().Post(url, p, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.PhoneInfo.PhoneNumber == "" {
		err = errors.New("获取用户手机号信息失败")
		return
	}
	return
}

type WxACodeParams struct {
	Path      string `json:"path"`
	Width     int    `json:"width,omitempty"`
	AutoColor bool   `json:"auto_color,omitempty"`
	LineCode  struct {
		R string `json:"r"`
		G string `json:"g"`
		B string `json:"b"`
	} `json:"line_color,omitempty"`
	IsHyaLine bool `json:"is_hyaline,omitempty"`
}

type WxACodeResp struct {
	WxCommonResp
	ContentType string      `json:"contentType"`
	Buffer      interface{} `json:"buffer"`
}

func (rec *api) GetWxACode(access string, params *WxACodeParams) (res *WxACodeResp, err error) {
	url := "https://api.weixin.qq.com/wxa/getwxacode?access_token=" + access
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
