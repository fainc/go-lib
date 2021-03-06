package default_middleware

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func MultiLang(r *ghttp.Request) {
	langArr := gstr.Explode(",", r.Header.Get("Accept-Language"))
	lang := "zh-CN" // 默认中文
	// languages:
	//  - zh-CN
	//  - en
	//  - jp
	languages, _ := g.Cfg().Get(r.Context(), "languages", garray.Array{}) // 从配置中读取

	langList := garray.NewStrArrayFrom(gconv.Strings(languages))
	if langArr != nil && langArr[0] != "" && langList.Contains(langArr[0]) {
		lang = langArr[0]
	}
	if r.Header.Get("Accept-Language") != "zh-CN" {
		r.SetCtx(gi18n.WithLanguage(r.Context(), lang))
	}
	r.SetCtxVar("API_VERSION", 1011)
	r.Middleware.Next()
}

func DefaultLangCN(r *ghttp.Request) {
	r.SetCtx(gi18n.WithLanguage(r.Context(), "zh-CN"))
	r.Middleware.Next()
}

func DefaultLangEN(r *ghttp.Request) {
	r.SetCtx(gi18n.WithLanguage(r.Context(), "en"))
	r.Middleware.Next()
}
