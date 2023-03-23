package md5_crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(s string, toUpper ...bool) string {
	h := md5.New()
	h.Write([]byte(s))
	if len(toUpper) == 1 && toUpper[0] {
		return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
