package wechat_sdk

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type utils struct{}

var utilsVar = utils{}

func Utils() *utils {
	return &utilsVar
}

// GetNonceStr 获取随机字符串
func (rec *utils) GetNonceStr() string {
	u1 := uuid.NewString()
	return strings.ToUpper(strings.ReplaceAll(u1, "-", ""))
}

// DownloadPathCheck 文件下载路径检查
func (rec *utils) DownloadPathCheck(path string) error {
	if path == "" {
		return errors.New("downloadPath不能为空")
	}
	lastStr := path[len(path)-1:]
	if lastStr != "/" {
		return errors.New("downloadPath应该以/结尾")
	}
	return nil
}

// HyaLineSuffix 根据IsHyaLine判断下载文件后缀
func (rec *utils) HyaLineSuffix(i bool) (suffix string) {
	suffix = ".jpeg"
	if i {
		suffix = ".png"
	}
	return
}
