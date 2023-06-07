package test

import (
	"fmt"
	"testing"

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
