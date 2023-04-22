package gm_crypto

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
)

func SM3Sum(data string, returnHex bool, salt ...string) (output string) {
	h := sm3.New()
	str := data
	if len(salt) >= 1 {
		str = str + salt[0]
	}
	h.Write([]byte(str))
	sum := h.Sum(nil)
	return formatRet(sum, returnHex)
}

func formatRet(sum []byte, returnHex bool) (output string) {
	if returnHex {
		return strings.ToUpper(hex.EncodeToString(sum))
	}
	return base64.StdEncoding.EncodeToString(sum)
}
func SM3FileSum(filePath string, returnHex bool) (output string, err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("SM3FileSum 读取文件失败")
		return
	}
	h := sm3.New()
	h.Write(f)
	sum := h.Sum(nil)
	output = formatRet(sum, returnHex)
	return
}
