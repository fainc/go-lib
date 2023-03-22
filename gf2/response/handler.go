package response

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
)

// HandlerResponse 默认数据返回中间件
func HandlerResponse(r *ghttp.Request) {
	var (
		ctx  = r.Context()
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)

	// openapi
	if gstr.Contains(r.RequestURI, "api.json") {
		return
	}

	// 已有err错误
	if err != nil {
		if code.Code() == 50 || code.Code() == 52 || code.Code() == 500 { // 服务器错误
			Json().InternalError(ctx, g.I18n().Translate(ctx, "InternalError"))
			return
		}
		if code.Code() == 401 { // 登录
			Json().UnAuthorizedError(ctx, code.Message(), code.Detail())
			return
		}
		if code.Code() == 402 { // 登录
			Json().DecryptError(ctx, code.Message(), code.Detail())
			return
		}
		if code.Code() == 403 { // 登录
			Json().SignatureError(ctx, code.Message(), code.Detail())
			return
		}
		if code.Code() == 404 { // 登录
			Json().NotFoundError(ctx, code.Message())
			return
		}
		Json().Error(ctx, err.Error(), code.Code(), code.Detail()) // 常规错误
		return
	}

	// 已退出程序流程
	if r.IsExited() {
		return
	}

	// 已有非错误自定义输出内容
	if gmeta.Get(res, "mime").String() == "custom" {
		return
	}

	// 已有异常响应状态码
	if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		switch r.Response.Status {
		case http.StatusNotFound: // 404
			Json().NotFoundError(ctx, g.I18n().Translate(ctx, "NotFound"))
			return
		case http.StatusUnauthorized: // 401
			return
		case http.StatusBadRequest: // 400
			return
		default:
			Json().InternalError(ctx, g.I18n().Translate(ctx, "InternalError"))
			return
		}
	}

	Json().Success(ctx, res)
}
