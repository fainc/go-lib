package wechat_sdk

type sat struct{}

var satVar = sat{}

func Sat() *sat {
	return &satVar
}

// Get  获取服务端接口调用凭据
func (rec *sat) Get(sdk *SdkClient) (token string, err error) {

	switch Cache().GetEngine() {
	case "redis":
		token, err = Cache().GetRedisCache("sat", sdk.AppId)
		if err != nil {
			return
		}
	default:
		token, err = Cache().GetMemoryCache("sat", sdk.AppId)
		if err != nil {
			return
		}
	}
	if token == "" {
		token, err = rec.Refresh(sdk)
		return
	}
	return
}

// Refresh  刷新接口调用凭据
func (rec *sat) Refresh(sdk *SdkClient) (token string, err error) {
	s, err := Api().GetSat(sdk)
	if err != nil {
		return
	}
	token = s.AccessToken
	switch Cache().GetEngine() {
	case "redis":
		err = Cache().SetRedisCache("sat", sdk.AppId, token, s.ExpiresIn)
		if err != nil {
			return
		}
	default:
		err = Cache().SetMemoryCache("sat", sdk.AppId, token, s.ExpiresIn)
		if err != nil {
			return
		}
	}
	return
}
