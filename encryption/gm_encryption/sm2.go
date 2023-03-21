package gm_encryption

import (
	"crypto/rand"
	"encoding/base64"
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
func SM2GenerateKeyPem(pwd string) (pri string, pub string, err error) {
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
func EncryptAsn1(pubPem string, data string) (cipherText string, err error) {
	pub, err := x509.ReadPublicKeyFromPem([]byte(pubPem))
	if err != nil {
		return
	}
	cipher, err := pub.EncryptAsn1([]byte(data), rand.Reader) // sm2加密
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func DecryptAsn1(priPem string, pwd string, data string) (plainText string, err error) {
	var password []byte
	if pwd != "" {
		password = []byte(pwd)
	}
	pri, err := x509.ReadPrivateKeyFromPem([]byte(priPem), password)
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.New("待解密数据处理失败")
		return
	}
	plain, err := pri.DecryptAsn1(d) // sm2解密
	if err != nil {
		err = errors.New("数据解密失败，请核对私钥证书是否正确")
		return
	}
	return string(plain), nil
}

func PrivateSign(priPem string, pwd string, data string) (signStr string, err error) {
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
	return base64.StdEncoding.EncodeToString(sign), nil
}

func PublicVerify(pubPem string, data string, sign string) (ok bool, err error) {
	pub, err := x509.ReadPublicKeyFromPem([]byte(pubPem))
	if err != nil {
		err = errors.New("加载私钥证书失败，请检查私钥证书和证书密码（若有）")
		return
	}
	sd, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		err = errors.New("签名数据处理失败")
		return
	}
	ok = pub.Verify([]byte(data), sd) // sm2签名
	return
}
