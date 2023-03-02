package wechat_sdk

type mpMessage struct{}

var mpMessageVar = mpMessage{}

func MpMessage() *mpMessage {
	return &mpMessageVar
}

func (rec *mpMessage) Method() (err error) {
	// code here
	return
}
