package default_controller

/**
控制系统停机，设置停机后通过中间件拦截请求（需要全局中间件注册）
*/
import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/fainc/go-lib/gf2/default_api"
)

var Shutdown = cShutdown{}

type cShutdown struct{}

var IsShutdown = false
var ShutdownMsg = ""
var ShutdownUnlockUri = ""

// ShutdownLock 系统停机锁定
func (c *cShutdown) ShutdownLock(ctx context.Context, req *default_api.ShutdownLockReq) (res *default_api.ShutdownRes, err error) {
	IsShutdown = true
	ShutdownMsg = "系统更新中，请稍后再试试"
	if req.Message != "" {
		ShutdownMsg = req.Message
	}
	uri := g.RequestFromCtx(ctx).RequestURI

	ShutdownUnlockUri = gstr.SubStr(uri, 0, gstr.Pos(uri, "/shutdown/lock")) + "/shutdown/unlock"
	return
}

// ShutdownUnlock 系统停机解锁
func (c *cShutdown) ShutdownUnlock(ctx context.Context, req *default_api.ShutdownUnlockReq) (res *default_api.ShutdownRes, err error) {
	IsShutdown = false
	ShutdownMsg = ""
	return
}
