package wechat_sdk

type miniProgramAnalyze struct{}

var miniProgramAnalyzeVar = miniProgramAnalyze{}

func MiniProgramAnalyze() *miniProgramAnalyze {
	return &miniProgramAnalyzeVar
}

func (rec *miniProgramAnalyze) Method() (err error) {
	// code here
	return
}
