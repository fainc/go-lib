package wechat_sdk

type miniProgramBroadcast struct{}

var miniProgramBroadcastVar = miniProgramBroadcast{}

func MiniProgramBroadcast() *miniProgramBroadcast {
	return &miniProgramBroadcastVar
}

func (rec *miniProgramBroadcast) Method() (err error) {
	// code here
	return
}
