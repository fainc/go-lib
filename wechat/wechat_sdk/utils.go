package wechat_sdk

import (
	"strings"

	"github.com/google/uuid"
)

type utils struct{}

var utilsVar = utils{}

func Utils() *utils {
	return &utilsVar
}

// GetNonceStr 获取随机字符串
func (rec *utils) GetNonceStr() string {
	u1 := uuid.NewString()
	return strings.ToUpper(strings.ReplaceAll(u1, "-", ""))
}
