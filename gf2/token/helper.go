package token

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/fainc/go-lib/jwt"
)

type helper struct{}

var helperVar = helper{}

func Helper() *helper {
	return &helperVar
}

type AuthJwtParams struct {
	Algo         string
	Aud          string
	WhiteTables  g.SliceStr
	Secret       string
	PublicKey    *ecdsa.PublicKey
	CryptoSecret string
}
type PublishJwtParams struct {
	Algo         string
	Aud          []string
	UserID       int64
	Secret       string
	PrivateKey   *ecdsa.PrivateKey
	CryptoAlgo   string
	CryptoSecret string
}

// PublishAuthToken 发布auth token(可直接使用或作为示例自行开发)
func (rec *helper) PublishAuthToken(p PublishJwtParams) (token, jti string, err error) {
	token, jti, err = jwt.Issuer(jwt.IssuerConf{
		JwtAlgo:      p.Algo,
		JwtPrivate:   p.PrivateKey,
		JwtSecret:    p.Secret,
		CryptoAlgo:   p.CryptoAlgo,
		CryptoSecret: p.CryptoSecret,
	}).Publish(&jwt.IssueParams{
		Subject:  "Auth",
		UserID:   gconv.String(p.UserID),
		Duration: 7 * 24 * time.Hour, // 授权24小时
		Audience: p.Aud,
		Ext:      "",
	})
	return
}

// StatelessAuth  用户无状态通用jwt验证和ctx写入(可直接使用或作为示例自行开发),通过err和catchErr判断拦截
func (rec *helper) StatelessAuth(r *ghttp.Request, p AuthJwtParams) (userID int64, aud string, catchErr bool, err error) {
	c, err := jwt.Parser(jwt.ParserConf{
		JwtAlgo:      p.Algo,
		JwtSecret:    p.Secret,
		CryptoSecret: p.CryptoSecret,
		JwtPublic:    p.PublicKey,
	}).Validate(jwt.ValidateParams{
		Token:    r.GetHeader("Authorization"),
		Subject:  "Auth", // 需要确保颁发时的subject = Auth才能通过验证，详见 Publish Auth 参数
		Audience: p.Aud,
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
	userID = gconv.Int64(c.UserID)
	aud = p.Aud
	r.SetCtxVar("TOKEN_UID", userID)
	r.SetCtxVar("TOKEN_JTI", c.ID)
	r.SetCtxVar("TOKEN_AUD", aud)
	return
}

type CtxJwtUser struct {
	ID      int64
	TokenID string
	Aud     string
}

// GetCtxUser 获取CTX用户ID信息
func (*helper) GetCtxUser(ctx context.Context) CtxJwtUser {
	r := g.RequestFromCtx(ctx)
	return CtxJwtUser{
		ID:      r.GetCtxVar("TOKEN_UID", 0).Int64(),
		TokenID: r.GetCtxVar("TOKEN_JTI", "").String(),
		Aud:     r.GetCtxVar("TOKEN_AUD", "").String(),
	}
}
