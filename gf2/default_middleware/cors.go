package default_middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// DefaultCors 内置CORS跨域中间件，建议注册
func DefaultCors(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
