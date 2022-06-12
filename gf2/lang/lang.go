package lang

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

func MiddlewareLang(r *ghttp.Request) {
	langArr := gstr.Explode(",", r.Header.Get("Accept-Language"))
	lang := "zh-CN" // 默认中文
	langList := garray.NewStrArrayFrom([]string{"zh-CN", "en", "fr"})
	if langArr != nil && langArr[0] != "" && langList.Contains(langArr[0]) {
		lang = langArr[0]
	}
	if r.Header.Get("Accept-Language") != "zh-CN" {
		r.SetCtx(gi18n.WithLanguage(r.Context(), lang))
	}
	r.SetCtxVar("API_VERSION", 1011)
	r.Middleware.Next()
}
