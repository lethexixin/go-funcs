package dysmsapi

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

type DySMS struct {
	Client *dysmsapi.Client
	opts   *Options
}

type Options struct {
	regionId        string
	accessKeyId     string
	accessKeySecret string
	scheme          string
	signName        string
	templateCode    string
}

type Option func(*Options)

func (d *DySMS) NewSMSClient(options ...Option) {
	opts := Options{}

	for _, o := range options {
		o(&opts)
	}

	client, err := dysmsapi.NewClientWithAccessKey(opts.regionId, opts.accessKeyId, opts.accessKeySecret)
	if err != nil {
		logger.Errorf("init dysms err: %s", err.Error())
	}

	d.Client = client
}

func (d *DySMS) NewSMSCodeRequest(phone string, code string) *dysmsapi.SendSmsRequest {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = d.opts.scheme
	request.PhoneNumbers = phone
	request.SignName = d.opts.signName
	request.TemplateCode = d.opts.templateCode
	request.TemplateParam = code
	return request
}

func (d *DySMS) SendSmsInfo(request *dysmsapi.SendSmsRequest) (err error) {
	_, err = d.Client.SendSms(request)
	if err != nil {
		return err
	}
	return nil
}
