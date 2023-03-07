package wechat_sdk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type aesCrypto struct{}

var aesCryptoVar = aesCrypto{}

func Aes() *aesCrypto {
	return &aesCryptoVar
}

func (rec *aesCrypto) pkcs5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// @brief:去除填充数据
func (rec *aesCrypto) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt AES加密
func (rec *aesCrypto) AesEncrypt(origData, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = rec.pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	cryptedStr := base64.StdEncoding.EncodeToString(crypted)
	return cryptedStr, nil
}

// AesDecrypt AES解密
func (rec *aesCrypto) AesDecrypt(str string, key []byte) ([]byte, error) {
	crypted, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = rec.pkcs5UnPadding(origData)
	return origData, nil
}
