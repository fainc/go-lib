package wechat_sdk

import (
	"errors"
)

type api struct{}

var apiVar = api{}

func Api() *api {
	return &apiVar
}

type SatRes struct {
	WxCommonRes
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetSat 小程序/公众号通用获取AccessToken方法
func (rec *api) GetSat(sdk *SdkClient) (res *SatRes, err error) {
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

type MpUserAccessTokenRes struct {
	WxCommonRes
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	OpenId         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotUser int    `json:"is_snapshotuser"`
	UnionId        string `json:"unionid"`
}

func (rec *api) GetMpUserAccessToken(sdk *SdkClient, code string) (res *MpUserAccessTokenRes, err error) {
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

type MpUserInfoRes struct {
	WxCommonRes
	OpenId     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	UnionId    string `json:"unionid"`
}

func (rec *api) GetMpUserInfo(accessToken string, openid string, lang string) (res *MpUserInfoRes, err error) {
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

type MpAccountUserInfoRes struct {
	WxCommonRes
	Subscribe      int    `json:"subscribe"`
	Openid         string `json:"openid"`
	Language       string `json:"language"`
	SubscribeTime  int    `json:"subscribe_time"`
	UnionId        string `json:"unionid"`
	Remark         string `json:"remark"`
	GroupId        int    `json:"groupid"`
	TagIdList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

func (rec *api) GetMpAccountUserInfo(accessToken string, openid string, lang string) (res *MpAccountUserInfoRes, err error) {
	if lang == "" {
		lang = "zh_CN"
	}
	url := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + accessToken + "&openid=" + openid + "&lang=" + lang
	err = Request().Get(url, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	return
}

type JsApiTicketRes struct {
	WxCommonRes
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

func (rec *api) GetJsApiTicket(accessToken string) (res *JsApiTicketRes, err error) {
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

type Code2SessionRes struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	WxCommonRes
}

func (rec *api) Code2Session(sdk *SdkClient, code string) (res *Code2SessionRes, err error) {
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

type PaidUnionIdParams struct {
	OpenId        string `json:"openid"`         // * 用户标识
	TransactionId string `json:"transaction_id"` // 可选，微信支付订单号
	MchId         string `json:"mch_id"`         // 可选，微信支付分配的商户号，和商户订单号配合使用
	OutTradeNo    string `json:"out_trade_no"`   // 可选，微信支付商户订单号，和商户号配合使用
}

type PaidUnionIdRes struct {
	UnionId string `json:"unionid"`
	WxCommonRes
}

func (rec *api) GetPaidUnionId(access string, p *PaidUnionIdParams) (res *PaidUnionIdRes, err error) {
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

// UserPhoneNumberRes 为减少Unmarshal错误的可能性，仅保留最小字段PurePhoneNumber、CountryCode，PhoneNumber和watermark忽略
type UserPhoneNumberRes struct {
	WxCommonRes
	PhoneInfo struct {
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

func (rec *api) GetUserPhoneNumber(access string, code string) (res *UserPhoneNumberRes, err error) {
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
	if res.PhoneInfo.PurePhoneNumber == "" {
		err = errors.New("获取用户手机号信息失败")
		return
	}
	return
}

type LineCodeParams struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}
type WxACodeParams struct {
	Path      string          `json:"path"`
	Width     int             `json:"width,omitempty"`
	AutoColor bool            `json:"auto_color,omitempty"`
	LineCode  *LineCodeParams `json:"line_color,omitempty"`
	IsHyaLine bool            `json:"is_hyaline,omitempty"`
}

func (rec *api) DownloadWxACode(access string, params *WxACodeParams, downloadPath string) (err error) {
	res := &WxCommonRes{}
	url := "https://api.weixin.qq.com/wxa/getwxacode?access_token=" + access
	err = Request().PostAndDownloadCode(url, params, &res, downloadPath)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	return
}

type WxACodeUnLimitParams struct {
	Scene      string          `json:"scene"`
	Page       string          `json:"page,omitempty"`
	CheckPath  bool            `json:"check_path"`
	EnvVersion string          `json:"env_version,omitempty"`
	Width      int             `json:"width,omitempty"`
	AutoColor  bool            `json:"auto_color,omitempty"`
	LineCode   *LineCodeParams `json:"line_color,omitempty"`
	IsHyaLine  bool            `json:"is_hyaline,omitempty"`
}

func (rec *api) DownloadWxACodeUnLimit(access string, params *WxACodeUnLimitParams, downloadPath string) (err error) {
	res := &WxCommonRes{}
	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + access
	err = Request().PostAndDownloadCode(url, params, &res, downloadPath)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	return
}

type WxAQrCodeParams struct {
	Path  string `json:"path"`
	Width int    `json:"width,omitempty"`
}

func (rec *api) DownloadWxAQrCode(access string, params *WxAQrCodeParams, downloadPath string) (err error) {
	res := &WxCommonRes{}
	url := "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=" + access
	err = Request().PostAndDownloadCode(url, params, &res, downloadPath)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	return
}

type JumpWxaParams struct {
	Path       string `json:"path"`
	Query      string `json:"query"`
	EnvVersion string `json:"env_version"`
}

type GenerateSchemeParams struct {
	JumpWxa        *JumpWxaParams `json:"jump_wxa,omitempty"`
	IsExpire       bool           `json:"is_expire,omitempty"`
	ExpireTime     int            `json:"expire_time,omitempty"`
	ExpireType     int            `json:"expire_type,omitempty"`
	ExpireInterval int            `json:"expire_interval,omitempty"`
}

type GenerateSchemeRes struct {
	WxCommonRes
	OpenLink string `json:"openLink"`
}

func (rec *api) GenerateScheme(access string, params *GenerateSchemeParams) (res *GenerateSchemeRes, err error) {
	url := "https://api.weixin.qq.com/wxa/generatescheme?access_token=" + access
	err = Request().Post(url, params, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.OpenLink == "" {
		err = errors.New("生成Scheme失败")
		return
	}
	return
}

type GenerateNFCSchemeParams struct {
	JumpWxa *JumpWxaParams `json:"jump_wxa,omitempty"`
	ModelId string         `json:"model_id"`
	Sn      string         `json:"sn,omitempty"`
}

func (rec *api) GenerateNFCScheme(access string, params *GenerateNFCSchemeParams) (res *GenerateSchemeRes, err error) {
	url := "https://api.weixin.qq.com/wxa/generatenfcscheme?access_token=" + access
	err = Request().Post(url, params, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.OpenLink == "" {
		err = errors.New("生成Scheme失败")
		return
	}
	return
}

type QuerySchemeParams struct {
	Scheme string `json:"scheme"`
}
type QuerySchemeRes struct {
	WxCommonRes
	SchemeInfo struct {
		AppId      string `json:"appid"`
		Path       string `json:"path"`
		Query      string `json:"query"`
		CreateTime int    `json:"create_time"`
		ExpireTime int    `json:"expire_time"`
		EnvVersion string `json:"env_version"`
	} `json:"scheme_info"`
	VisitOpenId string `json:"visit_openid"`
}

func (rec *api) QueryScheme(access string, params *QuerySchemeParams) (res *QuerySchemeRes, err error) {
	url := "https://api.weixin.qq.com/wxa/queryscheme?access_token=" + access
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

type CloudBaseParams struct {
	Env    string `json:"env"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
	Query  string `json:"query"`
}
type GenerateUrlLinkParams struct {
	Path           string           `json:"path,omitempty"`
	Query          string           `json:"query,omitempty"`
	IsExpire       bool             `json:"is_expire,omitempty"`
	ExpireType     int              `json:"expire_type,omitempty"`
	ExpireTime     int              `json:"expire_time,omitempty"`
	ExpireInterval int              `json:"expire_interval,omitempty"`
	EnvVersion     string           `json:"env_version,omitempty"`
	CloudBase      *CloudBaseParams `json:"cloud_base,omitempty"`
}
type GenerateUrlLinkRes struct {
	WxCommonRes
	UrlLink string `json:"url_link"`
}

func (rec *api) GenerateUrlLink(access string, params *GenerateUrlLinkParams) (res *GenerateUrlLinkRes, err error) {
	url := "https://api.weixin.qq.com/wxa/generate_urllink?access_token=" + access
	err = Request().Post(url, params, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.UrlLink == "" {
		err = errors.New("生成UrlLink失败")
		return
	}
	return
}

type QueryUrlLinkParams struct {
	UrlLink string `json:"url_link"`
}

type QueryUrlLinkRes struct {
	WxCommonRes
	UrlLinkInfo struct {
		AppId      string `json:"appid"`
		Path       string `json:"path"`
		Query      string `json:"query"`
		CreateTime int    `json:"create_time"`
		ExpireTime int    `json:"expire_time"`
		EnvVersion string `json:"env_version"`
	} `json:"url_link_info"`
	VisitOpenId string `json:"visit_openid"`
}

func (rec *api) QueryUrlLink(access string, params *QueryUrlLinkParams) (res *QueryUrlLinkRes, err error) {
	url := "https://api.weixin.qq.com/wxa/query_urllink?access_token=" + access
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

type GenerateShortLinkParams struct {
	PageUrl     string `json:"page_url"`
	PageTitle   string `json:"page_title,omitempty"`
	IsPermanent bool   `json:"is_permanent,omitempty"`
}
type GenerateShortLinkRes struct {
	WxCommonRes
	Link string `json:"link"`
}

func (rec *api) GenerateShortLink(access string, params *GenerateShortLinkParams) (res *GenerateShortLinkRes, err error) {
	url := "https://api.weixin.qq.com/wxa/genwxashortlink?access_token=" + access
	err = Request().Post(url, params, &res)
	if err != nil {
		return
	}
	if res.ErrCode != 0 {
		err = errors.New(res.ErrMsg)
		return
	}
	if res.Link == "" {
		err = errors.New("生成ShortLink失败")
		return
	}
	return
}

type MessageValue struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}
type MessageMiniProgramConf struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}
type MpTemplateMsg struct {
	AppId       string                   `json:"appid"`
	TemplateId  string                   `json:"template_id"`
	Url         string                   `json:"url"`
	MiniProgram *MessageMiniProgramConf  `json:"miniprogram"`
	Data        map[string]*MessageValue `json:"data"`
}

type UniformMessageParams struct {
	ToUser        string         `json:"touser"`
	MpTemplateMsg *MpTemplateMsg `json:"mp_template_msg"`
}

func (rec *api) SendUniformMessage(access string, params *UniformMessageParams) (res *WxCommonRes, err error) {
	url := "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send?access_token=" + access
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

type MiniProgramSubscribeMessageParams struct {
	ToUser           string                   `json:"touser"`
	TemplateId       string                   `json:"template_id"`
	Page             string                   `json:"page,omitempty"`
	MiniProgramState string                   `json:"miniprogram_state"`
	Lang             string                   `json:"lang"`
	Data             map[string]*MessageValue `json:"data"`
}

func (rec *api) SendSubscribeMessage(access string, params *MiniProgramSubscribeMessageParams) (res *WxCommonRes, err error) {
	url := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + access
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
