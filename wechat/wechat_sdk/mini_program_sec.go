package wechat_sdk

type miniProgramSec struct{}

var miniProgramSecVar = miniProgramSec{}

func MiniProgramSec() *miniProgramSec {
	return &miniProgramSecVar
}

func (rec *miniProgramSec) Method() (err error) {
	// code here
	return
}
