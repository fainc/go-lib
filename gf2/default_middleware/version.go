package default_middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gbuild"
)

// BindBuildVersion 系统版本（bin）中间件
func BindBuildVersion(r *ghttp.Request) {
	// 需要在build varMap 定义
	// buildVersion: xxx
	r.SetCtxVar("BUILD_VERSION", gbuild.Get("buildVersion", "UNKNOWN"))
	r.Middleware.Next()
}
