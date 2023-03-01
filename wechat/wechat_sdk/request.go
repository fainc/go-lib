package wechat_sdk

import (
	"encoding/json"
	"io"
	"net/http"
)

type request struct{}

var requestVar = request{}

func Request() *request {
	return &requestVar
}

type WxCommonResp struct {
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
