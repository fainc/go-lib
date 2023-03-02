package wechat_sdk

type mpAnalyze struct{}

var mpAnalyzeVar = mpAnalyze{}

func MpAnalyze() *mpAnalyze {
	return &mpAnalyzeVar
}

func (rec *mpAnalyze) Method() (err error) {
	// code here
	return
}
