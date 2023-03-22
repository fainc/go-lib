package gm_crypto

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
)

func SM3Sum(data string, salt string, toUpper bool) string {
	h := sm3.New()
	h.Write([]byte(data + salt))
	sum := h.Sum(nil)
	if toUpper {
		return strings.ToUpper(hex.EncodeToString(sum))
	}
	return hex.EncodeToString(sum)
}
func SM3FileSum(filePath string, toUpper bool) (ret string, err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("SM3FileSum 读取文件失败")
		return
	}
	h := sm3.New()
	h.Write(f)
	sum := h.Sum(nil)
	if toUpper {
		return strings.ToUpper(hex.EncodeToString(sum)), nil
	}
	return hex.EncodeToString(sum), nil
}
