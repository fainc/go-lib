package response

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
)

// MiddlewareHandlerResponse 默认数据返回中间件
func MiddlewareHandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()
	var (
		ctx = r.Context()
		err = r.GetError()
		res = r.GetHandlerResponse()
	)
	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 && gmeta.Get(res, "mime").String() == "custom" {
		return
	}
	if err != nil {
		_, _ = Json().Error(ctx, err.Error())
		return
	}
	_, _ = Json().Success(ctx, res)
}
