package default_middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/fainc/go-lib/gf2/default_controller"
	"github.com/fainc/go-lib/gf2/response"
)

// ShutdownLock 停机中间件
func ShutdownLock(r *ghttp.Request) {
	if default_controller.IsShutdown && r.RequestURI != default_controller.ShutdownUnlockUri {
		response.Json().Error(r.Context(), default_controller.ShutdownMsg, 401, nil)
		return
	}
	r.Middleware.Next()
}
