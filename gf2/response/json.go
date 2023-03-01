package response

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

type json struct{}

var jsonVar = json{}

func Json() *json {
	return &jsonVar
}

func (rec *json) Success(ctx context.Context, data interface{}) {
	rec.Writer(ctx, data, "Success", 200, 200, 0, nil)
}

func (rec *json) Error(ctx context.Context, message string, errCode int, ext interface{}) {
	rec.Writer(ctx, nil, message, 400, 400, errCode, ext)
}

func (rec *json) ServerError(ctx context.Context, message string) {
	rec.Writer(ctx, nil, message, 500, 500, 500, nil)
}

func (rec *json) Authorization(ctx context.Context, message string, ext interface{}) {
	rec.Writer(ctx, nil, message, 401, 401, 401, ext)
}

func (rec *json) NotFound(ctx context.Context, message string) {
	rec.Writer(ctx, nil, message, 404, 404, 404, nil)
}

type JsonFormat struct {
	Code         int         `json:"code"`         // 业务码，用于业务判断（首选），固定 200/400/401/404/500，与http状态码同步，通常取该值判断是否有错误需要处理
	ErrorCode    int         `json:"errorCode"`    // 错误码，用于业务判断（可选），-1/400(通用错误)/51(参数验证错误)/401/404/500/other，通常忽略该值，除非业务需要判断详细错误类型（例：交易场景，交易失败返回400业务码时，返回余额不足、账户冻结等详细错误码用于后续业务处理）
	Message      interface{} `json:"message"`      // 消息，业务码非400时固定输出 Success/Authorization/NotFound/InternalError中的一种，400时输出详细错误（可能含i18n转译，建议仅展示或记录信息，不可用于业务判断）
	Data         interface{} `json:"data"`         // 返回数据
	TraceId      interface{} `json:"traceId"`      // 请求唯一追踪ID,用于日志定位
	Server       string      `json:"server"`       // 服务端 host name md5 值
	Runtime      int64       `json:"runtime"`      // 当前任务执行耗时（ms）
	Lang         string      `json:"lang"`         // 当前上下文语言
	Ext          interface{} `json:"ext"`          // 拓展数据（可能含有多个错误详情或其他附加数据，例：复杂登录场景下的401错误返回具体登录地址）
	BuildVersion interface{} `json:"buildVersion"` // 当前程序运行版本号
	BootTime     int64       `json:"bootTime"`     // 系统启动时长（s）（使用时间戳差值而不是time.sub进行计算，time.sub受系统单调时钟策略影响，部分系统如mac休眠后单调时钟会停止，导致计算产生偏差）
	BootAt       string      `json:"bootAt"`       // 系统启动时间（应用启动时初始化）
}

// Writer 标准格式数据输出
func (rec *json) Writer(ctx context.Context, data interface{}, message string, status int, code int, errCode int, ext interface{}) {
	r := g.RequestFromCtx(ctx)
	r.Response.WriteStatus(status)
	r.Response.ClearBuffer()
	serverName, _ := os.Hostname()
	serverName, _ = gmd5.Encrypt(serverName)
	r.Response.WriteJson(JsonFormat{
		Code:         code,
		Message:      message,
		Data:         data,
		TraceId:      gctx.CtxId(ctx),
		Server:       strings.ToUpper(serverName),
		Ext:          ext,
		Lang:         gi18n.LanguageFromCtx(ctx),
		ErrorCode:    errCode,
		Runtime:      gtime.Now().TimestampMilli() - r.EnterTime,
		BuildVersion: r.GetCtxVar("BUILD_VERSION", "UNKNOWN"),
		BootTime:     GetBootTime(),
		BootAt:       GetBootAt(),
	})
	r.Response.Header().Set("Content-Type", "application/json;charset=utf-8")
	r.ExitAll()
}

// CustomWriter  自定义json输出
func (rec *json) CustomWriter(ctx context.Context, data interface{}, status ...int) {
	statusCode := 200
	if len(status) >= 1 && status[0] != 200 {
		statusCode = status[0]
	}
	r := g.RequestFromCtx(ctx)
	r.Response.WriteStatus(statusCode)
	r.Response.ClearBuffer()
	r.Response.WriteJson(data)
	r.Response.Header().Set("Content-Type", "application/json;charset=utf-8")
	r.ExitAll()
}
