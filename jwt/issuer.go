package jwt

import (
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fainc/go-lib/crypto/aes_crypto"
	"github.com/fainc/go-lib/crypto/gm_crypto"
	"github.com/fainc/go-lib/helper/str_helper"
)

type issuer struct {
	conf *IssuerConf
}

type IssuerConf struct {
	JwtAlgo      string            `dc:"* 指定JWT签名算法，支持 ES256(建议)，HS256(不建议)"`
	JwtPrivate   *ecdsa.PrivateKey `dc:"*(根据算法二选一) jwt签名私钥证书，根据签名算法选择，ES256应传私钥加签"`
	JwtSecret    string            `dc:"*(根据算法二选一) jwt签名密钥，根据签名算法选择，HS256应传不低于32位字符密钥加签"`
	CryptoAlgo   string            `dc:"可选，加密算法 支持AES和SM4(CBC模式)，不传则不进行加密"`
	CryptoSecret string            `dc:"可选，加密密钥，加密字段：UserID Ext"`
}

func Issuer(conf IssuerConf) *issuer {
	return &issuer{
		&conf,
	}
}

type IssueParams struct {
	Subject   string        `json:"subject"`          // * jwt主题键，如：UserAuth 用户验证 , Access 临时权限验证等
	UserID    string        `json:"userID"`           // * 用户编码
	Duration  time.Duration `json:"duration"`         // * 授权时长
	Audience  []string      `json:"audience"`         // 可选，授权作用域列表，验证时可判断授权是否在颁发列表内
	NotBefore time.Time     `json:"notBefore"`        // 可选，启用时间
	Ext       string        `json:"ext,omitempty"`    // 可选，额外用户信息，例如邮箱、昵称等，不建议存储用户敏感数据，如存储敏感数据请传加密密钥进行加密。
	JwtID     string        `json:"jwtID,omitempty"`  // 可选，自定义 jti，不传使用随机uuid
	Issuer    string        `json:"issuer,omitempty"` // 可选，签发者标记（可用于分布式签发端标记等
}

// check 签发基础信息核验
func (rec *issuer) check(params *IssueParams) (err error) {
	// 基础配置校验
	switch rec.conf.JwtAlgo { // 签名算法校验
	case AlgoHS256:
		if len(rec.conf.JwtSecret) < 32 {
			return errors.New("HS256签发算法 需要32位以上密钥")
		}
	case AlgoES256:
		if rec.conf.JwtPrivate == nil {
			return errors.New("ES256签发算法 私钥不能为空")
		}
	default:
		return errors.New("指定签发算法无效")
	}
	if rec.conf.CryptoAlgo != "" { // 加密算法密码强度校验
		switch rec.conf.CryptoAlgo {
		case AlgoAES:
			if len(rec.conf.CryptoSecret) != 32 {
				return errors.New("AES加密需要32位密钥(AES 256)")
			}
		case AlgoSM4:
			if len(rec.conf.CryptoSecret) != 16 {
				return errors.New("SM4加密需要16位密钥(SM4 128)")
			}
		default:
			return errors.New("不支持的加密算法")
		}
	}
	// 签发参数校验
	if params.UserID == "" || params.Duration <= 0 || params.Subject == "" {
		return errors.New("签发基础参数错误")
	}
	return
}
func (rec *issuer) defaultParams(params *IssueParams) {
	if params.JwtID == "" {
		params.JwtID = str_helper.UuidStr()
	}
	if params.NotBefore.IsZero() {
		params.NotBefore = time.Now()
	}
	// if params.Issuer == "" {
	// 	params.Issuer = "jwt.iss"
	// }
}

// encrypt 签发数据加密
func (rec *issuer) encrypt(params *IssueParams) (err error) {
	if rec.conf.CryptoAlgo == AlgoAES { // AES 256 CBC 加密
		if params.UserID != "" {
			if params.UserID, err = aes_crypto.EncryptCBC(params.UserID, rec.conf.CryptoSecret); err != nil {
				return
			}
		}
		if params.Ext != "" {
			if params.Ext, err = aes_crypto.EncryptCBC(params.Ext, rec.conf.CryptoSecret); err != nil {
				return
			}
		}
	}
	if rec.conf.CryptoAlgo == AlgoSM4 { // SM4 128 CBC 加密
		if params.UserID != "" {
			if params.UserID, err = gm_crypto.SM4Encrypt("CBC", rec.conf.CryptoSecret, params.UserID, false); err != nil {
				return
			}
		}
		if params.Ext != "" {
			if params.Ext, err = gm_crypto.SM4Encrypt("CBC", rec.conf.CryptoSecret, params.Ext, false); err != nil {
				return
			}
		}
	}
	return
}

// Publish 颁发token
func (rec *issuer) Publish(params *IssueParams) (token, jwtID string, err error) {
	// 基础信息校验
	if err = rec.check(params); err != nil {
		return
	}
	// 默认值处理
	rec.defaultParams(params)
	// 数据加密
	if rec.conf.CryptoAlgo != "" {
		if err = rec.encrypt(params); err != nil {
			return
		}
	}
	// 构造JWT
	claims := TokenClaims{
		params.UserID,
		params.Ext,
		rec.conf.CryptoAlgo,
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
	var method jwt.SigningMethod
	var secret interface{}
	if rec.conf.JwtAlgo == AlgoES256 { // ES256签发
		method = jwt.SigningMethodES256
		secret = rec.conf.JwtPrivate
	} else {
		method = jwt.SigningMethodHS256
		secret = []byte(rec.conf.JwtSecret)
	}
	t := jwt.NewWithClaims(method,
		claims)
	if token, err = t.SignedString(secret); err != nil {
		return
	}
	return token, params.JwtID, err
}
