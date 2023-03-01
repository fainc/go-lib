package wechat_payment_v2

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
)

type request struct{}

var requestVar = request{}

func Request() *request {
	return &requestVar
}

type commonResp struct {
	ReturnCode string `json:"return_code" xml:"return_code"`
	ReturnMsg  string `json:"return_msg" xml:"return_msg"`
}

// Send 请求方法
func (rec *request) Send(url string, params interface{}) (respBody []byte, err error) {
	x, err := xml.Marshal(params)
	if err != nil {
		return
	}
	body := strings.NewReader(string(x))
	hc := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	respBody, _ = io.ReadAll(resp.Body)
	common := &commonResp{}
	err = xml.Unmarshal(respBody, &common)
	if err != nil {
		return nil, nil
	}
	if common.ReturnCode != "SUCCESS" {
		return nil, errors.New(common.ReturnMsg)
	}
	return
}
