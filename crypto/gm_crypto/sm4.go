package gm_crypto

import (
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/tjfoc/gmsm/sm4"
)

func sm4Operate(key, data []byte, mode string, isEncrypt bool) (out []byte, err error) {
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
func SM4Encrypt(mode, key, data string, isHex bool) (outStr string, err error) {
	if data == "" {
		err = errors.New("不支持空内容加密")
		return
	}
	out, err := sm4Operate([]byte(key), []byte(data), mode, true)
	if err != nil {
		return
	}
	if isHex {
		return hex.EncodeToString(out), nil
	}
	return base64.StdEncoding.EncodeToString(out), nil
}
func SM4Decrypt(mode, key, data string, isHex bool) (outStr string, err error) {
	if data == "" {
		return "", nil
	}
	var db []byte
	if isHex {
		db, err = hex.DecodeString(data)
		if err != nil {
			err = errors.New("处理待解密数据失败")
			return
		}
	} else {
		db, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			err = errors.New("处理待解密数据失败")
			return
		}
	}
	out, err := sm4Operate([]byte(key), db, mode, false)
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
