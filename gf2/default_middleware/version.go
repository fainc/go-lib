package default_middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gbuild"
)

func ApiVersion(r *ghttp.Request) {
	// 需要在build varMap 定义
	// apiVersion: xxx
	r.SetCtxVar("API_VERSION", gbuild.Get("apiVersion"))
	r.Middleware.Next()
}
