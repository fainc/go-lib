package gm_encryption

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io/ioutil"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

func sm2generateKey() (key *sm2.PrivateKey, err error) {
	key, err = sm2.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	if !key.Curve.IsOnCurve(key.X, key.Y) {
		err = errors.New("生成密钥非SM2曲线")
		return
	}
	return
}

// SM2GenerateKeyPem 生成带格式的国密 SM2 密钥
func SM2GenerateKeyPem(pwd string) (pri, pub string, priHex, pubHex string, err error) {
	var password []byte
	if pwd != "" {
		password = []byte(pwd)
	}
	key, err := sm2generateKey()
	if err != nil {
		return
	}
	priByte, err := x509.WritePrivateKeyToPem(key, password) // 生成密钥文件
	if err != nil {
		return
	}
	pubKey, _ := key.Public().(*sm2.PublicKey)
	pubByte, err := x509.WritePublicKeyToPem(pubKey) // 生成公钥文件
	if err != nil {
		return
	}
	pri = string(priByte)
	pub = string(pubByte)

	// 解决前后端对接问题，hex密钥输出给JS端使用
	// 前后端对接注意JS 是否需要要给公钥和密文处理04标识
	// https://github.com/JuneAndGreen/sm-crypto/issues/42

	// 私钥 hex
	priHex = hex.EncodeToString(key.D.Bytes())
	// 公钥 hex（如需要04标识请自行添加）
	pubHex = hex.EncodeToString(pubKey.X.Bytes()) + hex.EncodeToString(pubKey.Y.Bytes())
	return
}

func SM2ReadPrivateKeyFromPem(priPem string, password string) (pri *sm2.PrivateKey, err error) {
	pri, err = x509.ReadPrivateKeyFromPem([]byte(priPem), []byte(password))
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	return
}
func SM2ReadPrivateKeyFromPath(filePath string, password string) (pri *sm2.PrivateKey, err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("读取密钥证书文件失败")
		return
	}
	pri, err = x509.ReadPrivateKeyFromPem(f, []byte(password))
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	return
}

func SM2ReadPublicKeyFromPem(pubPem string) (pub *sm2.PublicKey, err error) {
	pub, err = x509.ReadPublicKeyFromPem([]byte(pubPem))
	if err != nil {
		err = errors.New("加载公钥证书失败，请检查证书")
		return
	}
	return
}

func SM2ReadPublicKeyFromPath(filePath string) (pub *sm2.PublicKey, err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("读取密钥证书文件失败")
		return
	}
	pub, err = x509.ReadPublicKeyFromPem(f)
	if err != nil {
		err = errors.New("加载公钥证书失败，请检查证书")
		return
	}
	return
}
func EncryptAsn1(pubPem string, data string, mode string) (cipherText string, err error) {
	pub, err := x509.ReadPublicKeyFromPem([]byte(pubPem))
	if err != nil {
		return
	}
	cipher, err := pub.EncryptAsn1([]byte(data), rand.Reader) // sm2加密
	if err != nil {
		return
	}
	if mode == "hex" {
		return hex.EncodeToString(cipher), nil
	}
	return base64.StdEncoding.EncodeToString(cipher), nil

}

// Encrypt mode 0 C1C3C2 mode1 C1C2C3
func Encrypt(pubPem string, data string, outFormat string, mode int) (cipherText string, err error) {
	pub, err := x509.ReadPublicKeyFromPem([]byte(pubPem))
	if err != nil {
		return
	}
	cipher, err := sm2.Encrypt(pub, []byte(data), rand.Reader, mode)
	if err != nil {
		return
	}
	if outFormat == "hex" {
		return hex.EncodeToString(cipher), nil
	}
	return base64.StdEncoding.EncodeToString(cipher), nil

}
func DecryptAsn1(priPem string, pwd string, data string, mode string) (plainText string, err error) {
	var password []byte
	if pwd != "" {
		password = []byte(pwd)
	}
	pri, err := x509.ReadPrivateKeyFromPem([]byte(priPem), password)
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	var d []byte
	if mode != "hex" {
		d, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			err = errors.New("待解密数据处理失败")
			return
		}
	} else {
		d, err = hex.DecodeString(data)
		if err != nil {
			err = errors.New("待解密数据处理失败")
			return
		}
	}

	plain, err := pri.DecryptAsn1(d) // sm2解密
	if err != nil {
		err = errors.New("数据解密失败，请核对私钥证书是否正确")
		return
	}
	return string(plain), nil
}

func Decrypt(priPem string, pwd string, data string, inFormat string, mode int) (plainText string, err error) {
	var password []byte
	if pwd != "" {
		password = []byte(pwd)
	}
	pri, err := x509.ReadPrivateKeyFromPem([]byte(priPem), password)
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	var d []byte
	if inFormat != "hex" {
		d, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			err = errors.New("待解密数据处理失败")
			return
		}
	} else {
		d, err = hex.DecodeString(data)
		if err != nil {
			err = errors.New("待解密数据处理失败")
			return
		}
	}

	plain, err := sm2.Decrypt(pri, d, mode)
	if err != nil {
		err = errors.New("数据解密失败，请核对私钥证书是否正确")
		return
	}
	return string(plain), nil
}

func PrivateSign(priPem string, pwd string, data string, mode string) (signStr string, err error) {
	var password []byte
	if pwd != "" {
		password = []byte(pwd)
	}
	pri, err := x509.ReadPrivateKeyFromPem([]byte(priPem), password)
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	sign, err := pri.Sign(rand.Reader, []byte(data), nil) // sm2签名
	if err != nil {
		err = errors.New("签名失败，请检查私钥证书")
		return
	}
	if mode == "hex" {
		return hex.EncodeToString(sign), nil
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func PublicVerify(pubPem string, data string, sign string, mode string) (ok bool, err error) {
	pub, err := x509.ReadPublicKeyFromPem([]byte(pubPem))
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	var sd []byte
	if mode == "hex" {
		sd, err = hex.DecodeString(sign)
		if err != nil {
			err = errors.New("签名数据处理失败")
			return
		}
	} else {
		sd, err = base64.StdEncoding.DecodeString(sign)
		if err != nil {
			err = errors.New("签名数据处理失败")
			return
		}
	}
	ok = pub.Verify([]byte(data), sd) // sm2签名
	return
}
