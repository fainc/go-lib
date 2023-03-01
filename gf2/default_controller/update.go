package default_controller

/**
系统更新操作
更新前需手动shutdown系统，否则不予执行更新
*/

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/fainc/go-lib/gf2/default_api"
	"github.com/fainc/go-lib/gf2/response"
)

var Update = cUpdate{}

type cUpdate struct{}

// UpdateBin 通过上传bin文件更新并重启系统
func (c *cUpdate) UpdateBin(ctx context.Context, req *default_api.UpdateBinReq) (res *default_api.UpdateRes, err error) {
	res = &default_api.UpdateRes{}
	if !IsShutdown {
		err = errors.New("系统暂未停机，无法使用停机更新")
		return
	}
	pass := g.Cfg().MustGet(context.Background(), "server.managerPassword").String()
	if pass == "" {
		err = errors.New("当前系统未启用更新密码")
		return
	}
	if pass != req.Password {
		err = errors.New("更新密码错误")
		return
	}
	req.Bin.Filename = gtime.Now().TimestampStr() + "_" + grand.S(4)
	_, err = req.Bin.Save("./", false)
	if err != nil {
		return
	}
	err = ghttp.RestartAllServer(ctx, req.Bin.Filename)
	if err != nil {
		return
	}
	return
}

// UpdateShowUpload 内置上传bin文件更新页面
func (c *cUpdate) UpdateShowUpload(ctx context.Context, req *default_api.UpdateShowUploadReq) (res *response.CustomRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Write(`
    <html>
    <head>
        <title>上传更新</title>
    </head>
        <body>
            <form enctype="multipart/form-data" action="./bin-update" method="post">
                <input type="file" name="bin" required/>
                <input type="password" name="password" placeholder="请输入更新密码" required />
                <input type="submit" value="确定更新" />
            </form>
        </body>
    </html>
    `)
	return
}
