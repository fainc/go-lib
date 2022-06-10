package response

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func AuthorizedError(ext string) error {
	return gerror.NewCode(gcode.CodeNotAuthorized, ext)
}
