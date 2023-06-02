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
// Deprecated: 请使用新版接口 jwt.Validator 验证Jwt
func (*jwtHelper) Parse(params ParseParams) (userID int, scope, jwtId string, claims jwt.MapClaims, err error) {
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
	jtiStr := claims["jti"]
	if jtiStr == nil {
		err = errors.New("jti invalid")
		return
	}
	return gconv.Int(userIdStr), gconv.String(scopeStr), gconv.String(jtiStr), claims, nil
}

type GenerateParams struct {
	UserId    int64         // * 非0用户ID
	Scope     string        // * 授权scope标志
	Duration  time.Duration // * 授权时长
	Secret    string        // * jwt及加密密钥
	Id        string        // * 唯一标识，为空时使用随机uuid xxxx-xxxx-xxxx-xxx
	NotBefore *time.Time    // 生效时间 nil时使用now
	Ext       string        // 附加数据
}

// Generate 生成jwt
// Deprecated: 请使用新版接口 jwt.Helper().Issue 发行token
func (*jwtHelper) Generate(params GenerateParams) (tokenString, jti string, err error) {
	if params.UserId == 0 || params.Scope == "" || params.Duration == 0 || params.Secret == "" {
		err = errors.New("generate jwt params invalid")
		return
	}
	if params.Id == "" {
		params.Id = str_helper.UuidStr()
	}
	type MyCustomClaims struct {
		UserId int64  `json:"userId"`
		Scope  string `json:"scope"`
		Ext    string `json:"ext"`
		jwt.RegisteredClaims
	}
	nbf := time.Now()
	if params.NotBefore != nil {
		nbf = *params.NotBefore
	}
	claims := MyCustomClaims{
		params.UserId,
		params.Scope,
		params.Ext,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(nbf),
			Issuer:    "jwtHelper",
			ID:        params.Id,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(params.Secret))
	return
}

type StandardAuthParams struct {
	Scopes      g.SliceStr
	WhiteTables g.SliceStr
	Secret      string
}

// StandardAuth 通用jwt验证和ctx写入(可直接使用或作为示例自行开发),通过err和catchErr判断拦截
func (rec *jwtHelper) StandardAuth(r *ghttp.Request, p StandardAuthParams) (userID int, scope string, catchErr bool, err error) {
	userID, scope, jwtID, _, err := rec.Parse(ParseParams{
		Token:  r.GetHeader("Authorization"),
		Scopes: p.Scopes,
		Secret: p.Secret,
	})
	if err != nil {
		catchErr = true
		if len(p.WhiteTables) != 0 {
			whiteTable := garray.NewStrArrayFrom(p.WhiteTables)
			if whiteTable.ContainsI(r.URL.Path) {
				catchErr = false
			}
		}
		return
	}
	r.SetCtxVar("JWT_USER_ID", userID)
	r.SetCtxVar("JWT_USER_TOKEN_ID", jwtID)
	r.SetCtxVar("JWT_USER_SCOPE", scope)
	return
}

type CtxJwtUser struct {
	ID      int
	TokenID string
	Scope   string
}

// GetCtxUser 获取CTX用户ID信息
func (*jwtHelper) GetCtxUser(ctx context.Context) CtxJwtUser {
	r := g.RequestFromCtx(ctx)
	return CtxJwtUser{
		ID:      r.GetCtxVar("JWT_USER_ID", 0).Int(),
		TokenID: r.GetCtxVar("JWT_USER_TOKEN_ID", "").String(),
		Scope:   r.GetCtxVar("JWT_USER_SCOPE", "").String(),
	}
}
