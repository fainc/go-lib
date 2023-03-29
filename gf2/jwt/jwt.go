package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"

	"github.com/fainc/go-lib/helper/str_helper"
)

var Helper = jwtHelper{}

type jwtHelper struct{}

type ParseParams struct {
	Token  string     // * jwt字符串
	Scopes g.SliceStr // * jwt scope可用范围
	Secret string     // * jwt密钥
}

// Parse jwt解析
func (*jwtHelper) Parse(params ParseParams) (userId int, scope, id string, claims jwt.MapClaims, err error) {
	if params.Secret == "" {
		err = errors.New("jwt secret invalid")
		return
	}
	if params.Token == "" {
		err = errors.New("authorization invalid")
		return
	}
	tokenMap := strings.Split(params.Token, "Bearer ")
	if len(tokenMap) != 2 {
		err = errors.New("bearer invalid")
		return
	}
	tokenString := tokenMap[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(params.Secret), nil
	})
	if err != nil {
		err = errors.New(err.Error())
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		err = errors.New(err.Error())
		return
	}
	userIdStr := claims["userId"]
	if userIdStr == nil {
		err = errors.New("signature user key invalid")
		return
	}
	scopeStr := claims["scope"]
	scopes := garray.NewStrArrayFrom(params.Scopes)
	if scopeStr == nil || !scopes.ContainsI(gconv.String(scopeStr)) {
		err = errors.New("scope invalid")
		return
	}
	idStr := claims["jti"]
	if idStr == nil {
		err = errors.New("jti invalid")
		return
	}
	return gconv.Int(userIdStr), gconv.String(scopeStr), gconv.String(idStr), claims, nil
}

type GenerateParams struct {
	UserId    int           // * 非0用户ID
	Scope     string        // * 授权scope标志
	Duration  time.Duration // * 授权时长
	Secret    string        // * jwt及加密密钥
	Id        string        // * 唯一标识，为空时使用随机uuid xxxx-xxxx-xxxx-xxx
	NotBefore *time.Time    // * 生效时间 nil时使用now
}

// Generate 生成jwt
func (*jwtHelper) Generate(params GenerateParams) (string, error) {
	if params.UserId == 0 || params.Scope == "" || params.Duration == 0 || params.Secret == "" {
		return "", errors.New("generate jwt params invalid")
	}
	if params.Id == "" {
		params.Id = str_helper.UuidStr()
	}
	type MyCustomClaims struct {
		UserId int    `json:"userId"`
		Scope  string `json:"scope"`
		jwt.RegisteredClaims
	}
	nbf := time.Now()
	if params.NotBefore != nil {
		nbf = *params.NotBefore
	}
	claims := MyCustomClaims{
		params.UserId,
		params.Scope,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(nbf),
			Issuer:    "jwtHelper",
			ID:        params.Id,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(params.Secret))
	return tokenString, nil
}

// StandardAuth 通用jwt验证和ctx写入(可直接使用或作为示例自行开发)
func (rec *jwtHelper) StandardAuth(r *ghttp.Request, scopes g.SliceStr, secret string) (userId int, scope, id string, claims jwt.MapClaims, err error) {
	userId, scope, id, claims, err = rec.Parse(ParseParams{
		Token:  r.GetHeader("Authorization"),
		Scopes: scopes,
		Secret: secret,
	})
	if err != nil {
		scope = "UNKNOWN"
		return
	}
	r.SetCtxVar("JWT_USER_ID", userId)
	r.SetCtxVar("JWT_SCOPE", scope)
	return
}

type jwtUser struct {
	ID    int
	SCOPE string
}

// GetCtxUser 获取CTX用户信息
func (*jwtHelper) GetCtxUser(ctx context.Context) jwtUser {
	r := g.RequestFromCtx(ctx)
	return jwtUser{
		ID:    gconv.Int(r.GetCtxVar("JWT_USER_ID", 0)),
		SCOPE: gconv.String(r.GetCtxVar("JWT_SCOPE", "UNKNOWN")),
	}
}
