package wechat_sdk

type miniProgramMessage struct{}

var miniProgramMessageVar = miniProgramMessage{}

func MiniProgramMessage() *miniProgramMessage {
	return &miniProgramMessageVar
}

func (rec *miniProgramMessage) Method() (err error) {
	// code here
	return
}
