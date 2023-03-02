package wechat_sdk

type mpCus struct{}

var mpCusVar = mpCus{}

func MpCus() *mpCus {
	return &mpCusVar
}

func (rec *mpCus) Method() (err error) {
	// code here
	return
}
