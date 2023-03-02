package wechat_sdk

import (
	"fmt"
)

type miniProgram struct{}

var miniProgramVar = miniProgram{}

func MiniProgram() *miniProgram {
	return &miniProgramVar
}

// Code2Session 小程序登录，获取session和openId、unionId(开放平台关联才有该值)
func (rec *miniProgram) Code2Session(sdk *SdkClient, code string) (res *Code2SessionResp, err error) {
	res, err = Api().Code2Session(sdk, code)
	return
}

type GetPaidUnionIdParams struct {
	OpenId        string `json:"openid"`         // * 用户标识
	TransactionId string `json:"transaction_id"` // 可选，微信支付订单号
	MchId         string `json:"mch_id"`         // 可选，微信支付分配的商户号，和商户订单号配合使用
	OutTradeNo    string `json:"out_trade_no"`   // 可选，微信支付商户订单号，和商户号配合使用
}

// GetPaidUnionId 支付后获取unionId，调用前需要用户完成支付，且在支付后的五分钟内有效。
// 使用微信支付订单号（transaction_id）和微信支付商户订单号和微信支付商户号（out_trade_no 及 mch_id），二选一
func (rec *miniProgram) GetPaidUnionId(sdk *SdkClient, p *GetPaidUnionIdParams) (res *PaidUnionIdResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GetPaidUnionId(token, p)
	return
}

// GetUserPhoneNumber 获取小程序用户手机号
func (rec *miniProgram) GetUserPhoneNumber(sdk *SdkClient, code string) (res *UserPhoneNumberResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GetUserPhoneNumber(token, code)
	return
}

// DownloadWxACode 获取小程序码
//  doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getQRCode.html
func (rec *miniProgram) DownloadWxACode(sdk *SdkClient, params *WxACodeParams, downloadPath string) (path string, err error) {
	if err = Utils().DownloadPathCheck(downloadPath); err != nil {
		return
	}
	suffix := Utils().HyaLineSuffix(params.IsHyaLine)
	path = downloadPath + "WxACode_" + Utils().GetNonceStr() + suffix
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	err = Api().DownloadWxACode(token, params, path)
	return
}

// DownloadWxACodeUnLimit 获取不限制的小程序码
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getUnlimitedQRCode.html
func (rec *miniProgram) DownloadWxACodeUnLimit(sdk *SdkClient, params *WxACodeUnLimitParams, downloadPath string) (path string, err error) {
	if err = Utils().DownloadPathCheck(downloadPath); err != nil {
		return
	}
	if !params.CheckPath {
		fmt.Println("DownloadWxACodeUnLimit 警告：checkPath为 false 时允许小程序未发布或者 page 不存在， 但page 有数量上限（60000个）请勿滥用。")
	}
	suffix := Utils().HyaLineSuffix(params.IsHyaLine)
	path = downloadPath + "WxACodeUL_" + Utils().GetNonceStr() + suffix
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	err = Api().DownloadWxACodeUnLimit(token, params, path)
	return
}

// DownloadWxAQrCode 获取小程序二维码
//  doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/createQRCode.html
func (rec *miniProgram) DownloadWxAQrCode(sdk *SdkClient, params *WxAQrCodeParams, downloadPath string) (path string, err error) {
	if err = Utils().DownloadPathCheck(downloadPath); err != nil {
		return
	}
	path = downloadPath + "WxAQrCode_" + Utils().GetNonceStr() + ".jpeg"
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	err = Api().DownloadWxAQrCode(token, params, path)
	return
}

// GenerateScheme 生成小程序 Scheme URL
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/url-scheme/generateScheme.html
func (rec *miniProgram) GenerateScheme(sdk *SdkClient, params *GenerateSchemeParams) (res *GenerateSchemeResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GenerateScheme(token, params)
	return
}

// GenerateNFCScheme 生成小程序NFC Scheme URL
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/url-scheme/generateNFCScheme.html
func (rec *miniProgram) GenerateNFCScheme(sdk *SdkClient, params *GenerateNFCSchemeParams) (res *GenerateSchemeResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GenerateNFCScheme(token, params)
	return
}

// QueryScheme 查询Scheme
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/url-scheme/queryScheme.html
func (rec *miniProgram) QueryScheme(sdk *SdkClient, params *QuerySchemeParams) (res *QuerySchemeResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().QueryScheme(token, params)
	return
}

// GenerateUrlLink 生成UrlLink
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/url-link/generateUrlLink.html
func (rec *miniProgram) GenerateUrlLink(sdk *SdkClient, params *GenerateUrlLinkParams) (res *GenerateUrlLinkResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GenerateUrlLink(token, params)
	return
}

// QueryUrlLink 查询 UrlLink
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/url-link/queryUrlLink.html
func (rec *miniProgram) QueryUrlLink(sdk *SdkClient, params *QueryUrlLinkParams) (res *QueryUrlLinkResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().QueryUrlLink(token, params)
	return
}

// GenerateShortLink 生成小程序短链接 适用于微信内拉起小程序
// doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/short-link/generateShortLink.html
func (rec *miniProgram) GenerateShortLink(sdk *SdkClient, params *GenerateShortLinkParams) (res *GenerateShortLinkResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GenerateShortLink(token, params)
	return
}
