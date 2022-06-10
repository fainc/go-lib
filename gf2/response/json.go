package response

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

type json struct{}

var jsonVar = json{}

func Json() *json {
	return &jsonVar
}

func (rec *json) Success(ctx context.Context, data interface{}) {
	rec.Writer(ctx, data, "success", 200, 200, nil)
}

func (rec *json) Error(ctx context.Context, message interface{}) {
	rec.Writer(ctx, nil, message, 400, 400, nil)
}

func (rec *json) ServerError(ctx context.Context, message interface{}) {
	rec.Writer(ctx, nil, message, 500, 500, nil)
}

func (rec *json) Authorization(ctx context.Context, message interface{}) {
	rec.Writer(ctx, nil, message, 401, 401, nil)
}

func (rec *json) NotFound(ctx context.Context, message interface{}) {
	rec.Writer(ctx, nil, message, 404, 404, nil)
}

type JsonFormat struct {
	Code    int         `json:"code"`    // 业务码
	Message interface{} `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 返回数据
	TraceId interface{} `json:"traceId"` // 请求唯一追踪ID
	Server  string      `json:"server"`  // 服务端 host name md5 值
	Runtime int64       `json:"runtime"` // 执行时长（ms）
	Ext     interface{} `json:"ext"`     // 拓展数据（自定义调用Writer）
}

// Writer 数据输出
func (rec *json) Writer(ctx context.Context, data interface{}, message interface{}, status int, code int, ext interface{}) {
	r := g.RequestFromCtx(ctx) // 从Ctx中获取Request对象
	r.Response.WriteStatus(status)
	r.Response.ClearBuffer()
	r.Response.Header().Set("Content-Type", "application/json;charset=utf-8")
	serverName, _ := os.Hostname()
	serverName, _ = gmd5.Encrypt(serverName)
	_ = r.Response.WriteXml(JsonFormat{
		Code:    code,
		Message: message,
		Data:    data,
		TraceId: gctx.CtxId(ctx),
		Server:  strings.ToUpper(serverName),
		Ext:     ext,
		Runtime: gtime.Now().TimestampMilli() - r.EnterTime,
	})
	return
}
