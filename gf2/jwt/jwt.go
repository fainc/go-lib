package jwt

import (
	"context"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/fainc/go-lib/jwt"
)

var Helper = jwtHelper{}

type jwtHelper struct{}

type StandardAuthParams struct {
	Conf        jwt.ParserConf
	Params      jwt.ValidateParams
	WhiteTables g.SliceStr
}

// StandardAuth 通用jwt验证和ctx写入(可直接使用或作为示例自行开发),通过err和catchErr判断拦截
func (rec *jwtHelper) StandardAuth(r *ghttp.Request, p StandardAuthParams) (userID int64, aud string, catchErr bool, err error) {
	c, err := jwt.Parser(p.Conf).Validate(p.Params)
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
	aud = p.Params.Audience
	r.SetCtxVar("JWT_USER_ID", userID)
	r.SetCtxVar("JWT_USER_JTI", c.ID)
	r.SetCtxVar("JWT_USER_AUD", aud)
	return
}

type CtxJwtUser struct {
	ID      int64
	TokenID string
	Aud     string
}

// GetCtxUser 获取CTX用户ID信息
func (*jwtHelper) GetCtxUser(ctx context.Context) CtxJwtUser {
	r := g.RequestFromCtx(ctx)
	return CtxJwtUser{
		ID:      r.GetCtxVar("JWT_USER_ID", 0).Int64(),
		TokenID: r.GetCtxVar("JWT_USER_TOKEN_ID", "").String(),
		Aud:     r.GetCtxVar("JWT_USER_AUD", "").String(),
	}
}
