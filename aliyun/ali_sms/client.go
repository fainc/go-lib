package ali_sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysms "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var SmsClient = smsClientService{}

type smsClientService struct{}

func (s *smsClientService) NewClient(accessKeyId string, accessKeySecret string, endPoint string) (client *dysms.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	if endPoint == "" {
		endPoint = "dysmsapi.aliyuncs.com"
	}
	config.Endpoint = tea.String(endPoint)
	client, err = dysms.NewClient(config)
	return
}
