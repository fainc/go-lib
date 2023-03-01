package default_api

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ShutdownRes struct {
}
type ShutdownLockReq struct {
	g.Meta  `path:"/shutdown/lock" tags:"DefaultController" method:"get" summary:"停机维护设置"`
	Message string `json:"message" v:"required" dc:"停机信息"`
}
type ShutdownUnlockReq struct {
	g.Meta `path:"/shutdown/unlock" tags:"DefaultController" method:"get" summary:"取消停机维护"`
}
