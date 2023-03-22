package response

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func InternalError(ctx context.Context, ext interface{}) error {
	return gerror.NewCode(gcode.New(500, g.I18n().Translate(ctx, "InternalError"), ext))
}
func UnAuthorizedError(ctx context.Context, ext interface{}) error {
	return gerror.NewCode(gcode.New(401, g.I18n().Translate(ctx, "UnAuthorized"), ext))
}
func DecryptError(ctx context.Context, ext interface{}) error {
	return gerror.NewCode(gcode.New(402, g.I18n().Translate(ctx, "DecryptFailed"), ext))
}
func SignatureError(ctx context.Context, ext interface{}) error {
	return gerror.NewCode(gcode.New(403, g.I18n().Translate(ctx, "SignatureError"), ext))
}
func NotFoundError(ctx context.Context, ext interface{}) error {
	return gerror.NewCode(gcode.New(404, g.I18n().Translate(ctx, "NotFound"), ext))
}

// StandardError 通用错误返回，支持i18n 如需自定义错误码ENUM对应相关错误，请自行参考实现
func StandardError(ctx context.Context, code int, message string, trans bool, ext interface{}) error {
	if trans {
		message = g.I18n().Translate(ctx, message)
	}
	if code == 401 {
		return UnAuthorizedError(ctx, ext)
	}
	return gerror.NewCode(gcode.New(code, message, ext))
}
