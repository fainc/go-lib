package wechat_sdk

type miniProgramCus struct{}

var miniProgramCusVar = miniProgramCus{}

func MiniProgramCus() *miniProgramCus {
	return &miniProgramCusVar
}

func (rec *miniProgramCus) Method() (err error) {
	// code here
	return
}
