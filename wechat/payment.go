package wechat

type payment struct{}

var paymentVar = payment{}

func Payment() *payment {
	return &paymentVar
}

func (rec *payment) Method() (err error) {
	// code here
	return
}
