package wechat_sdk

type client struct{}

var clientVar = client{}

func Client() *client {
	return &clientVar
}

type SdkClient struct {
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}
