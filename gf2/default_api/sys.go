package default_api

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetStatusReq struct {
	g.Meta   `path:"/sys/status" tags:"DefaultController" method:"get" summary:"获取系统信息"`
	Password string `json:"password" v:"required"`
}
