package md5_crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(s string, toUpper bool, salt ...string) string {
	h := md5.New()
	if len(salt) >= 1 {
		s = s + salt[0]
	}
	h.Write([]byte(s))
	if toUpper {
		return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
