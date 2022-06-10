package response

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
)

// MiddlewareHandlerResponse 默认数据返回中间件
func MiddlewareHandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()
	var (
		ctx  = r.Context()
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	// api.json 不作处理
	if r.RequestURI == "/api.json" {
		return
	}
	// 已有自定义输出内容，不作处理
	if r.Response.BufferLength() > 0 && gmeta.Get(res, "mime").String() == "custom" {
		return
	}
	if err != nil {
		if code == gcode.CodeNil || code == gcode.CodeInternalError { // 服务器错误
			Json().ServerError(ctx, g.I18n().Translate(ctx, "InternalError"))
			return
		}
		if code == gcode.CodeNotAuthorized { // 登录
			Json().Authorization(ctx, g.I18n().Translate(ctx, "Authorization"), err.Error())
			return
		}
		Json().Error(ctx, err.Error()) // 常规错误
		return
	}
	if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		switch r.Response.Status {
		case http.StatusNotFound:
			Json().NotFound(ctx, g.I18n().Translate(ctx, "NotFound"))
			return
		case http.StatusForbidden:
			Json().Authorization(ctx, g.I18n().Translate(ctx, "Authorization"), nil)
			return
		default:
			Json().ServerError(ctx, g.I18n().Translate(ctx, "InternalError"))
			return
		}
	}
	Json().Success(ctx, res)
}
