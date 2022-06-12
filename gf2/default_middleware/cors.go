package default_middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func DefaultCors(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
