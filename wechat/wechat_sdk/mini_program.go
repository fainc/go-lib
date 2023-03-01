package wechat_sdk

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

// GetWxACode 获取小程序码
// params doc : https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getQRCode.html
func (rec *miniProgram) GetWxACode(sdk *SdkClient, params *WxACodeParams) (res *WxACodeResp, err error) {
	token, err := Sat().Get(sdk)
	if err != nil {
		return
	}
	res, err = Api().GetWxACode(token, params)
	return
}
