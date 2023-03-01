package wechat_sdk

// User Access Token 用户调用凭证
type uat struct{}

var uatVar = uat{}

func Uat() *uat {
	return &uatVar
}

func (rec *uat) Method() (err error) {
	// code here
	return
}
