package str_helper

import (
	"strings"

	"github.com/google/uuid"
)

// NonceStr 获取随机字符串
func NonceStr() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
}

// UuidStr 获取UUID
func UuidStr() string {
	return strings.ToUpper(uuid.NewString())
}
