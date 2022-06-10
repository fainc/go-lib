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
	// 已有自定义输出内容，不作处理
	if r.Response.BufferLength() > 0 && gmeta.Get(res, "mime").String() == "custom" {
		return
	}
	if err != nil { // 有错误信息
		Json().Error(ctx, err.Error())
		return
	}
	Json().Success(ctx, res)
}
