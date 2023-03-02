package wechat_sdk

type miniProgramHardware struct{}

var miniProgramHardwareVar = miniProgramHardware{}

func MiniProgramHardware() *miniProgramHardware {
	return &miniProgramHardwareVar
}

func (rec *miniProgramHardware) Method() (err error) {
	// code here
	return
}
