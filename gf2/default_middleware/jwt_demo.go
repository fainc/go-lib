package default_middleware

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/fainc/go-lib/gf2/jwt"
	"github.com/fainc/go-lib/gf2/response"
)

// JwtUserAuth this is a jwt middleware demo
func JwtUserAuth(r *ghttp.Request) {
	inWhiteTables := garray.NewStrArrayFrom(g.SliceStr{"/account/login"}).ContainsI(r.RequestURI)
	userId, scope, claims, err := jwt.Helper.StandardAuth(r, g.SliceStr{
		"scope",
	}, "secret")
	if err != nil && !inWhiteTables {
		response.Json().Authorization(r.Context(), err.Error(), nil)
		return // return 后下个中间件不执行
	}
	g.Dump(claims)
	g.Dump(scope)
	g.Dump(userId)
	g.Dump(inWhiteTables)
	r.Middleware.Next()
}
