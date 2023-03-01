package wechat_sdk

import (
	"bytes"
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
