package gm_crypto_v1

import (
	"encoding/base64"
	"errors"

	"github.com/tjfoc/gmsm/sm4"
)

func sm4Operate(key []byte, data []byte, mode string, isEncrypt bool) (out []byte, err error) {
	switch mode {
	case "ECB":
		out, err = sm4.Sm4Ecb(key, data, isEncrypt)
	case "CBC":
		out, err = sm4.Sm4Cbc(key, data, isEncrypt)
	case "CFB":
		out, err = sm4.Sm4CFB(key, data, isEncrypt)
	case "OFB":
		out, err = sm4.Sm4OFB(key, data, isEncrypt)
	default:
		err = errors.New("不支持的加解密模式：" + mode)
	}
	return
}
func SM4Encrypt(mode string, key string, data string) (outStr string, err error) {
	kb := []byte(key)
	db := []byte(data)
	out, err := sm4Operate(kb, db, mode, true)
	if err != nil {
		return
	}
	outStr = base64.StdEncoding.EncodeToString(out)
	return
}
func SM4Decrypt(mode string, key string, data string) (outStr string, err error) {
	kb := []byte(key)
	db, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.New("处理待解密数据失败")
		return
	}
	out, err := sm4Operate(kb, db, mode, false)
	if err != nil {
		return
	}
	if out == nil {
		err = errors.New("数据解密失败，请核实密钥")
		return
	}
	outStr = string(out)
	return
}
