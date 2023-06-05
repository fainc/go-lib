package aes_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// EncryptCBC AES CBC 加密
func EncryptCBC(text, key string, ivOption ...string) (string, error) {
	if text == "" {
		return "", errors.New("不支持空内容加密")
	}
	ivStr := "0000000000000000"
	if ivOption != nil || len(ivOption) >= 1 {
		ivStr = ivOption[0]
	}
	iv := []byte(ivStr)
	plaintext := []byte(text)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return "", errors.New("IV length must equal block size")
	}
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(encrypted, plaintext)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptCBC(text, key string, ivOption ...string) (string, error) {
	if text == "" {
		return "", nil
	}
	ivStr := "0000000000000000"
	if ivOption != nil || len(ivOption) >= 1 {
		ivStr = ivOption[0]
	}
	iv := []byte(ivStr)
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return "", errors.New("IV length must equal block size")
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = PKCS7UnPadding(origData)
	if origData == nil {
		return "", errors.New("AES解密错误，请核对密钥和密文")
	}
	return string(origData), nil
}
