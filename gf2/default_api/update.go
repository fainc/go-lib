package default_api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UpdateShowUploadReq struct {
	g.Meta `path:"/update/show-upload-bin" tags:"UpdateController" method:"get" summary:"更新系统"`
}
type UpdateBinReq struct {
	g.Meta   `path:"/update/bin-update" tags:"UpdateController" method:"post" summary:"更新系统"`
	Bin      *ghttp.UploadFile `p:"bin" json:"bin" v:"required" dc:"bin文件，file格式"`
	Password string            `json:"password" v:"required" dc:"更新密码"`
}
type UpdateRes struct {
	Version string `json:"version"`
}
