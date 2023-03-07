package wechat_sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type remoteCredentials struct{}

var remoteCredentialsVar = remoteCredentials{}

func RemoteCredentials() *remoteCredentials {
	return &remoteCredentialsVar
}

type getSatRes struct {
	remoteCredentialsCommonResp
	Data *struct {
		Sat string `json:"sat"`
	} `json:"data"`
}
type getJatRes struct {
	remoteCredentialsCommonResp
	Data *struct {
		Jat string `json:"jat"`
	} `json:"data"`
}

// GetSat 获取凭据中心 Server Access Token（远程需要实现缓存、刷新维护等功能）
func (rec *remoteCredentials) GetSat(sdk *SdkClient) (token string, err error) {
	params := make(map[string]string)
	params["appId"] = sdk.AppId
	params["secret"] = sdk.Secret
	res := &getSatRes{}
	err = rec.request("/credentials/sat/get", params, res)
	if err != nil {
		return
	}
	if res.Code != 200 {
		err = errors.New(res.Message)
		return
	}
	if res.Data == nil || res.Data.Sat == "" {
		err = errors.New("获取统一凭据中心Sat失败")
		return
	}
	return res.Data.Sat, nil
}

// GetJat  获取凭据中心 Js Api Ticket（远程需要实现缓存、刷新维护等功能）
func (rec *remoteCredentials) GetJat(sdk *SdkClient) (ticket string, err error) {
	params := make(map[string]string)
	params["appId"] = sdk.AppId
	params["secret"] = sdk.Secret
	res := &getJatRes{}
	err = rec.request("/credentials/jat/get", params, res)
	if err != nil {
		return
	}
	if res.Code != 200 {
		err = errors.New(res.Message)
		return
	}
	if res.Data == nil || res.Data.Jat == "" {
		err = errors.New("获取统一凭据中心Jat失败")
		return
	}
	return res.Data.Jat, nil
}

type remoteCredentialsCommonResp struct {
	Code      int    `json:"code"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

// request 请求方法，微信密钥和微信appId通过请求加密传输，远程凭据中心只负责缓存和维护，不保存密钥,credentialsAppId是凭据中心调用Id
func (rec *remoteCredentials) request(uri string, data map[string]string, res interface{}) (err error) {
	c, err := Cache().GetRemoteCredentialsClient()
	if err != nil {
		return
	}
	dataJson, _ := json.Marshal(data)
	dataAes, err := Aes().AesEncrypt(dataJson, []byte(c.Secret))
	if err != nil {
		return
	}
	params := make(map[string]string)
	params["appId"] = c.AppId // 统一凭据中心服务端调用凭证
	params["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	params["data"] = dataAes
	params["sign"] = rec.sign(params, c.Secret)

	paramsJson, _ := json.Marshal(params)
	body := bytes.NewBuffer(paramsJson)
	hc := &http.Client{}
	req, err := http.NewRequest("POST", c.Host+uri, body)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, res)
	return
}

// sign 签名方法，服务端自行实现请参考签名方法构造
func (rec *remoteCredentials) sign(params map[string]string, key string) (s string) {
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 对键进行排序
	var str string
	for _, v := range keys {
		if v != "sign" {
			str = str + v + "=" + params[v] + "&"
		}
	}
	str = str + "key=" + key
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
