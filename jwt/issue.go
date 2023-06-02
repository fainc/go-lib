package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fainc/go-lib/crypto/aes_crypto"
	"github.com/fainc/go-lib/helper/str_helper"
)

type IssueParams struct {
	Subject      string        `json:"subject"`                // * jwt主题键，如：Auth 用户验证 , Access 临时权限验证等
	UserID       string        `json:"userID"`                 // * 用户编码
	Audience     []string      `json:"audience"`               // * 授权作用域列表
	Duration     time.Duration `json:"duration"`               // * 授权时长
	JwtSecret    string        `json:"jwtSecret"`              // * jwt密钥，根据签名算法选择密钥，ES256应传私钥加签，HS256应传不低于16位字符密钥加签
	NotBefore    time.Time     `json:"notBefore"`              // 可选，启用时间
	Ext          string        `json:"ext,omitempty"`          // 可选，额外用户信息，例如邮箱、昵称等，不建议存储用户敏感数据，如存储敏感数据请传加密密钥进行加密。
	JwtAlgo      string        `json:"jwtAlg,omitempty"`       // 可选，自定义JWT签名算法，默认ES256(建议)，可选HS256(不建议)
	JwtID        string        `json:"jwtID,omitempty"`        // 可选，自定义 jti，不传使用随机uuid
	Issuer       string        `json:"issuer,omitempty"`       // 可选，签发者标记（可用于分布式签发端标记等），不传默认jwt.helper
	CryptoSecret string        `json:"cryptoSecret,omitempty"` // 可选，加密密钥，加密字段：UserID Ext，不传则不加密
	CryptoAlgo   string        `json:"cryptoAlgo,omitempty"`   // 可选，加密算法 目前仅支持AES(CBC模式)
	Redis        bool          `json:"redis"`                  // 可选，是否启用redis维护登陆态，启用redis后过期时间将由redis维护，不写入jwt内
}

func (rec *helper) GenKey() (pubStr, priStr string, err error) {
	generateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}
	derText, err := x509.MarshalECPrivateKey(generateKey)
	if err != nil {
		return
	}
	pubDerText, err := x509.MarshalPKIXPublicKey(&generateKey.PublicKey)
	if err != nil {
		return
	}
	priBlock := pem.Block{
		Type:  "ECDSA Private Key",
		Bytes: derText,
	}
	pri := pem.EncodeToMemory(&priBlock)
	pubBlock := pem.Block{
		Type:  "ECDSA Public Key",
		Bytes: pubDerText,
	}
	pub := pem.EncodeToMemory(&pubBlock)
	if pub == nil || pri == nil {
		err = errors.New("gen key failed")
		return
	}
	pubStr = string(pub)
	priStr = string(pri)
	return
}

// issueBase 签发基础信息核验
func (rec *helper) issueBase(params *IssueParams) (err error) {
	if params.UserID == "" || params.Duration == 0 || params.JwtSecret == "" || params.Subject == "" || params.Audience == nil {
		err = errors.New("generate jwt params invalid")
		return
	}
	// 默认值处理
	if params.JwtID == "" {
		params.JwtID = str_helper.UuidStr()
	}
	if params.NotBefore.IsZero() {
		params.NotBefore = time.Now()
	}
	if params.Issuer == "" {
		params.Issuer = "jwt.helper"
	}
	return
}

// issueEncrypt 签发数据加密
func (rec *helper) issueEncrypt(params *IssueParams) (err error) {
	if len(params.CryptoSecret) != 32 {
		err = errors.New("JWT加密密钥必须为32位(AES 256)，不支持16或24位短密钥")
		return
	}
	if params.UserID, err = aes_crypto.EncryptCBC(params.UserID, params.CryptoSecret); err != nil {
		return
	}
	if params.Ext != "" {
		if params.Ext, err = aes_crypto.EncryptCBC(params.Ext, params.CryptoSecret); err != nil {
			return
		}
	}
	return
}

// Issue 颁发token
func (rec *helper) Issue(params *IssueParams) (token, jwtID string, err error) {
	// 基础信息校验
	if err = rec.issueBase(params); err != nil {
		return
	}
	// 数据加密
	if params.CryptoAlgo == "AES" && params.CryptoSecret != "" {
		if err = rec.issueEncrypt(params); err != nil {
			return
		}
	}
	// 构造JWT详细信息
	type myCustomClaims struct {
		UserID interface{} `json:"userID"`
		Ext    string      `json:"ext"`
		jwt.RegisteredClaims
	}
	claims := myCustomClaims{
		params.UserID,
		params.Ext,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(params.NotBefore),
			Issuer:    params.Issuer,
			Subject:   params.Subject,
			ID:        params.JwtID,
			Audience:  params.Audience,
		},
	}
	// 签发
	var t *jwt.Token
	if params.JwtAlgo != "HS256" {
		t = jwt.NewWithClaims(jwt.SigningMethodES256,
			claims)
		var pri *ecdsa.PrivateKey
		if pri, err = jwt.ParseECPrivateKeyFromPEM([]byte(params.JwtSecret)); err != nil {
			err = errors.New("JWT签名私钥无法解析")
			return
		}
		token, err = t.SignedString(pri)
		return token, params.JwtID, err
	}
	if len(params.JwtSecret) < 32 {
		return "", "", errors.New("JWT签发密钥必须>=32位，不支持短密钥")
	}
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)
	token, err = t.SignedString([]byte(params.JwtSecret))
	return token, params.JwtID, err
}
