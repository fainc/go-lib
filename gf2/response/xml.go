package response

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type xml struct{}

var xmlVar = xml{}

func Xml() *xml {
	return &xmlVar
}

func (rec *xml) CustomWriter(ctx context.Context, data interface{}, status ...int) {
	statusCode := 200
	if len(status) >= 1 && status[0] != 200 {
		statusCode = status[0]
	}
	r := g.RequestFromCtx(ctx)
	r.Response.WriteStatus(statusCode)
	r.Response.ClearBuffer()
	r.Response.WriteXml(data, "xml")
	r.Response.Header().Set("Content-Type", "application/xml;charset=utf-8")
	r.ExitAll()
}
