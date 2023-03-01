package wechat_sdk

type miniProgram struct{}

var miniProgramVar = miniProgram{}

func MiniProgram() *miniProgram {
	return &miniProgramVar
}

func (rec *miniProgram) Method() (err error) {
	// code here
	return
}
