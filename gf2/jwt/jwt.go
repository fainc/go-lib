package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
)

var Helper = jwtHelper{}

type jwtHelper struct{}

type ParseParams struct {
	Token  string     // * jwt字符串
	Scopes g.SliceStr // * jwt scope可用范围
	Secret string     // * jwt密钥
}

// Parse jwt解析
func (*jwtHelper) Parse(params ParseParams) (int, string, error) {
	if params.Secret == "" {
		return 0, "", errors.New("jwt secret invalid")
	}
	if params.Token == "" {
		return 0, "", errors.New("authorization invalid")
	}
	tokenMap := strings.Split(params.Token, "Bearer ")
	if len(tokenMap) != 2 {
		return 0, "", errors.New("bearer invalid")
	}
	tokenString := tokenMap[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(params.Secret), nil
	})
	if err != nil {
		return 0, "", errors.New(err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New(err.Error())
	}
	uuid := claims["uuid"]
	if uuid == nil {
		return 0, "", errors.New("signature user key invalid")
	}
	scope := claims["scope"]
	scopes := garray.NewStrArrayFrom(params.Scopes)
	if scope == nil || !scopes.ContainsI(gconv.String(scope)) {
		return 0, "", errors.New("scope invalid")
	}
	return gconv.Int(uuid), gconv.String(scope), nil
}

type GenerateParams struct {
	Uuid     int           // * 非0用户ID
	Scope    string        // * 授权scope标志
	Duration time.Duration // * 授权时长
	Secret   string        // * jwt及加密密钥
}

// Generate 生成jwt
func (*jwtHelper) Generate(params GenerateParams) (string, error) {
	if params.Uuid == 0 || params.Scope == "" || params.Duration == 0 || params.Secret == "" {
		return "", errors.New("generate jwt params invalid")
	}

	type MyCustomClaims struct {
		Uuid  int    `json:"uuid"`
		Scope string `json:"scope"`
		jwt.RegisteredClaims
	}
	claims := MyCustomClaims{
		params.Uuid,
		params.Scope,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jwtHelper",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(params.Secret))
	return tokenString, nil
}

// StandardAuth 通用jwt验证和ctx写入(可直接使用或作为示例自行开发)
func (rec *jwtHelper) StandardAuth(r *ghttp.Request, scopes g.SliceStr, whiteTables g.SliceStr, secret string) (int, string, error) {
	whiteTable := garray.NewStrArrayFrom(whiteTables)
	uuid, scopeKey, err := rec.Parse(ParseParams{
		Token:  r.GetHeader("Authorization"),
		Scopes: scopes,
		Secret: secret,
	})
	if err != nil {
		if !whiteTable.ContainsI(r.RequestURI) {
			return 0, "UNKNOWN", err
		}
		r.SetCtxVar("UUID", 0)
		r.SetCtxVar("SCOPE", "UNKNOWN")
		return 0, "UNKNOWN", nil
	}
	r.SetCtxVar("UUID", uuid)
	r.SetCtxVar("SCOPE", scopeKey)
	return uuid, scopeKey, nil
}

type user struct {
	UUID  int
	SCOPE string
}

// GetUser 获取当前用户信息
func (*jwtHelper) GetUser(r *ghttp.Request) *user {
	UUID := r.GetCtxVar("UUID", 0)
	SCOPE := r.GetCtxVar("SCOPE", "UNKNOWN")
	return &user{
		UUID:  gconv.Int(UUID),
		SCOPE: gconv.String(SCOPE),
	}
}
