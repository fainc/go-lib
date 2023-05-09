package wechat_sdk

import (
	"net/url"
)

type mp struct {
	sdk *SdkClient
}

func Mp(AppId, Secret string) *mp {
	sdk, _ := Client().New(SdkClient{
		AppId:  AppId,
		Secret: Secret,
	})
	return &mp{sdk}
}

type AuthorizationParams struct {
	RedirectUri string `json:"redirectUri"` // * 重定向地址
	Scope       string `json:"scope"`       // * 默认snsapi_base；应用授权作用域，snsapi_base （不弹出授权页面，直接跳转，只能获取用户openid），snsapi_userinfo （弹出授权页面，可通过 openid 拿到昵称、性别、所在地。并且， 即使在未关注的情况下，只要用户授权，也能获取其信息 ）
	State       string `json:"state"`       // 可选 重定向后会带上 state 参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
}

// GetAuthorizationUrl 获取微信授权登录URL
func (rec *mp) GetAuthorizationUrl(params *AuthorizationParams) string {
	if params.Scope == "" {
		params.Scope = "snsapi_base"
	}
	params.RedirectUri = url.QueryEscape(params.RedirectUri)
	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + rec.sdk.AppId + "&redirect_uri=" + params.RedirectUri + "&response_type=code&scope=" + params.Scope + "&state=" + params.State + "#wechat_redirect"
}

// Code2AccessToken 微信公众号授权登录 获取OpenId和AccessToken
func (rec *mp) Code2AccessToken(code string) (res *MpUserAccessTokenRes, err error) {
	res, err = Api().GetMpUserAccessToken(rec.sdk, code)
	if err != nil {
		return
	}
	return
}

// GetUserInfo 获取用户信息（使用用户accessToken + openId）
func (rec *mp) GetUserInfo(userAccessToken string, openId string, lang string) (res *MpUserInfoRes, err error) {
	res, err = Api().GetMpUserInfo(userAccessToken, openId, lang)
	return
}

// GetAccountUserInfo 获取公众号用户信息（使用公众号accessToken + openId），注意：本接口不会输出头像昵称，未关注公众号拉取不到任何信息，主要用途：用于判断用户是否关注公众号
func (rec *mp) GetAccountUserInfo(openId string, lang string) (res *MpAccountUserInfoRes, err error) {
	token, err := Sat(rec.sdk).Get()
	if err != nil {
		return
	}
	res, err = Api().GetMpAccountUserInfo(token, openId, lang)
	return
}

// GetUserSubscribe 获取用户是否关注公众号（GetAccountUserInfo 的简化版）
func (rec *mp) GetUserSubscribe(openId string) (subscribe int, err error) {
	res, err := rec.GetAccountUserInfo(openId, "")
	if err != nil {
		return
	}
	subscribe = res.Subscribe
	return
}

// GetQrCode 获取公众号二维码
// https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html
func (rec *mp) GetQrCode(params *MpQrCodeParams, downloadPath string) (res *MpQrCodeRes, err error) {
	token, err := Sat(rec.sdk).Get()
	if err != nil {
		return
	}
	res, err = MpApi().GetQrCode(token, params)
	if err != nil {
		return
	}
	if downloadPath != "" {
		err = Utils().DownloadPathCheck(downloadPath)
		if err != nil {
			return
		}
		err = Request().GetDownload("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket="+url.QueryEscape(res.Ticket), downloadPath)
		return
	}
	return
}
