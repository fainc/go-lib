package response

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func AuthorizedError(ctx context.Context, ext interface{}) error {
	return gerror.NewCode(gcode.New(401, g.I18n().Translate(ctx, "Authorization"), ext))
}

func StandardError(ctx context.Context, code int, message string, trans bool, ext interface{}) error {
	if trans {
		message = g.I18n().Translate(ctx, message)
	}
	if code == 401 {
		return AuthorizedError(ctx, ext)
	}
	return gerror.NewCode(gcode.New(code, message, ext))
}
