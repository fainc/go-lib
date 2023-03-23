package gm_crypto

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
)

func SM3Sum(data string, salt ...string) (hex, bs64 string) {
	h := sm3.New()
	str := data
	if len(salt) >= 1 {
		str = str + salt[0]
	}
	h.Write([]byte(str))
	sum := h.Sum(nil)
	return formatRet(sum)
}

func formatRet(sum []byte) (hexStr, bs64 string) {
	return strings.ToUpper(hex.EncodeToString(sum)), base64.StdEncoding.EncodeToString(sum)
}
func SM3FileSum(filePath string) (hex, bs64 string, err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("SM3FileSum 读取文件失败")
		return
	}
	h := sm3.New()
	h.Write(f)
	sum := h.Sum(nil)
	hex, bs64 = formatRet(sum)
	return
}
