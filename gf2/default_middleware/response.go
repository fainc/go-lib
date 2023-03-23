package default_middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/fainc/go-lib/gf2/response"
)

func HandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()
	response.HandlerResponse(r)
}
func HandlerEncryptResponse(r *ghttp.Request) {
	r.Middleware.Next()
	response.HandlerResponse(r)
}
