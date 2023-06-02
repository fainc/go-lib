package jwt

type helper struct{}

var helperVar = helper{}

func Helper() *helper {
	return &helperVar
}

func (rec *helper) Method() (err error) {
	// code here
	return
}
