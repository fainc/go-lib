package wechat_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type request struct{}

var requestVar = request{}

func Request() *request {
	return &requestVar
}

type WxCommonRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (rec *request) Get(url string, res interface{}) (err error) {
	hc := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return
	}
	return
}

func (rec *request) Post(url string, data interface{}, res interface{}) (err error) {
	j, err := json.Marshal(data)
	if err != nil {
		return
	}
	hc := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return
	}
	return
}

func (rec *request) PostAndDownloadCode(url string, data interface{}, res interface{}, downloadPath string) (err error) {
	j, err := json.Marshal(data)
	if err != nil {
		return
	}
	hc := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		err = os.WriteFile(downloadPath, respBody, 0600)
		if err != nil {
			fmt.Println("PostAndDownloadCode 文件写入错误：" + err.Error())
			err = errors.New("文件保存失败，请检查保存路径及读写权限")
			return
		}
		return
	}
	return
}
func (rec *request) GetDownload(url string, downloadPath string) (err error) {
	hc := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("远程请求状态码错误")
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	err = os.WriteFile(downloadPath, respBody, 0600)
	if err != nil {
		fmt.Println("GetAndDownload 文件写入错误：" + err.Error())
		err = errors.New("文件保存失败，请检查保存路径及读写权限")
		return
	}
	return
}
