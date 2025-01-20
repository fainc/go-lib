package test

import (
	"fmt"
	"testing"

	"github.com/fainc/go-lib/wechat/wechat_payment_v1"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/fainc/go-lib/wechat/wechat_sdk"
)

func TestName(t *testing.T) {
	number, err := wechat_sdk.MiniProgram("123", "123").GetUserPhoneNumber("123")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(number)
}

func TestWxPayRefund(t *testing.T) {
	refund, err := wechat_payment_v1.Payment().Refund(&wechat_payment_v1.RefundRequestParams{
		OutTradeNo:  "WX25011709534799701",
		OutRefundNo: "WXR2501170953479970",
		TotalFee:    "5900",
		RefundFee:   "5900",
	}, GetWxPayClient())
	if err != nil {
		g.Dump(err)
		return
	}
	g.Dump(refund)
}
