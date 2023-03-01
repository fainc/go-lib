package wechat_sdk

import (
	"net/url"
)

type mp struct{}

var mpVar = mp{}

func Mp() *mp {
	return &mpVar
}

type AuthorizationParams struct {
	RedirectUri string `json:"redirectUri"` // * 重定向地址
	Scope       string `json:"scope"`       // * 默认snsapi_base；应用授权作用域，snsapi_base （不弹出授权页面，直接跳转，只能获取用户openid），snsapi_userinfo （弹出授权页面，可通过 openid 拿到昵称、性别、所在地。并且， 即使在未关注的情况下，只要用户授权，也能获取其信息 ）
	State       string `json:"state"`       // 可选 重定向后会带上 state 参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
}

// GetAuthorizationUrl 获取微信授权登录URL
func (rec *mp) GetAuthorizationUrl(sdk *SdkClient, params *AuthorizationParams) string {
	if params.Scope == "" {
		params.Scope = "snsapi_base"
	}
	params.RedirectUri = url.QueryEscape(params.RedirectUri)
	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + sdk.AppId + "&redirect_uri=" + params.RedirectUri + "&response_type=code&scope=" + params.Scope + "&state=" + params.State + "#wechat_redirect"
}

// Code2AccessToken 微信公众号授权登录 获取OpenId和AccessToken
func (rec *mp) Code2AccessToken(sdk *SdkClient, code string) (res *MpUserAccessTokenResp, err error) {
	res, err = Api().GetMpUserAccessToken(sdk, code)
	if err != nil {
		return
	}
	// todo uat缓存
	return
}
func (rec *mp) GetUserInfo(accessToken string, openId string, lang string) (res *MpUserInfoResp, err error) {
	res, err = Api().GetMpUserInfo(accessToken, openId, lang)
	return
}
