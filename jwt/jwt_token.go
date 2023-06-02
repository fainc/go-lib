package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/fainc/go-lib/helper/str_helper"
)

type jwtToken struct {
	Redis bool // 是否启用redis维护登陆态，启用redis后过期时间将由redis维护，不写入jwt内
}

func JwtToken(redis bool) *jwtToken {
	return &jwtToken{Redis: redis}
}

type IssueParams struct {
	UserID       int64         `json:"userID,omitempty"`       // * 用户编码
	Scope        string        `json:"scope,omitempty"`        // * 授权作用域
	Duration     time.Duration `json:"duration,omitempty"`     // * 授权时长
	JwtSecret    string        `json:"jwtSecret,omitempty"`    // * jwt密钥
	NotBefore    *time.Time    `json:"notBefore"`              // 可选，启用时间
	Ext          string        `json:"ext,omitempty"`          // 可选，额外用户信息，例如邮箱、昵称等，不建议存储用户敏感数据，如存储敏感数据请传加密密钥进行加密。
	JwtID        string        `json:"jwtID,omitempty"`        // 可选，自定义 jti，不传使用随机uuid
	CryptoSecret string        `json:"cryptoSecret,omitempty"` // 可选，加密密钥，加密字段：UserID Scope Ext，不传则不加密
	CryptoAlgo   string        `json:"cryptoAlgo,omitempty"`   // 可选，加密算法 支持 SM4、AES
	Redis        bool          `json:"redis"`                  // 可选，是否启用redis维护登陆态，启用redis后过期时间将由redis维护，不写入jwt内
}

// Issue 颁发token
func (rec *jwtToken) Issue(params IssueParams) (tokenString, jwtID string, err error) {
	if params.UserID == 0 || params.Scope == "" || params.Duration == 0 || params.JwtSecret == "" {
		err = errors.New("generate jwt params invalid")
		return
	}
	if params.JwtID == "" {
		params.JwtID = str_helper.UuidStr()
	}
	jwtID = params.JwtID
	type MyCustomClaims struct {
		UserID interface{} `json:"userID"`
		Scope  string      `json:"scope"`
		Ext    string      `json:"ext"`
		jwt.RegisteredClaims
	}
	nbf := time.Now()
	if params.NotBefore != nil {
		nbf = *params.NotBefore
	}
	claims := MyCustomClaims{
		params.UserID,
		params.Scope,
		params.Ext,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(nbf),
			Issuer:    "jwtHelper",
			ID:        params.JwtID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(params.JwtSecret))
	return
}
