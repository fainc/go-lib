package response

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func InternalError(ctx context.Context, ext interface{}) error {
	return Error(500, g.I18n().Translate(ctx, "InternalError"), ext)
}
func UnAuthorizedError(ctx context.Context, ext interface{}) error {
	return Error(401, g.I18n().Translate(ctx, "UnAuthorized"), ext)
}
func DecryptError(ctx context.Context, ext interface{}) error {
	return Error(402, g.I18n().Translate(ctx, "DecryptFailed"), ext)
}
func SignatureError(ctx context.Context, ext interface{}) error {
	return Error(403, g.I18n().Translate(ctx, "SignatureError"), ext)
}
func NotFoundError(ctx context.Context, ext interface{}) error {
	return Error(404, g.I18n().Translate(ctx, "NotFound"), ext)
}

// I18nError 返回一个标准错误信息，适用i18n 如需自定义错误码ENUM对应相关错误，请自行参考实现
func I18nError(ctx context.Context, code int, message string, ext interface{}) error {
	message = g.I18n().Translate(ctx, message)
	return Error(code, message, ext)
}

// Error 返回一个标准错误信息  code 建议使用400通用错误码（或自定义四位错误码与HTTP状态码区分开），一般不使用公共的401,402,403,404错误码，除非明确输出的是相关内容
func Error(code int, message string, ext interface{}) error {
	return gerror.NewCode(gcode.New(code, message, ext))
}
