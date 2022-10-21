package emails

import (
	"errors"
	"fmt"
)

import (
	"gopkg.in/gomail.v2"
)

type EMail struct {
	opts *Options
}

type Options struct {
	email    string // 发送者email
	password string // 发送者密码
	host     string // 发送者email对应的host,比如 smtphm.qiye.163.com
	port     int    // host对应的port,比如25
}

type Option func(*Options)

const (
	DefaultEmail    = "demo@163.com"
	DefaultPassword = "123456"
	DefaultHost     = "smtphm.qiye.163.com"
	DefaultPort     = 25
)

func Email(email string) Option {
	return func(o *Options) {
		o.email = email
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.password = password
	}
}

func Host(host string) Option {
	return func(o *Options) {
		o.host = host
	}
}

func Port(port int) Option {
	return func(o *Options) {
		o.port = port
	}
}

func (e *EMail) Init(options ...Option) {
	opts := Options{
		email:    DefaultEmail,
		password: DefaultPassword,
		host:     DefaultHost,
		port:     DefaultPort,
	}

	for _, o := range options {
		o(&opts)
	}

	e.opts = &opts
}

// SendEmail
// 发送邮件
// receiver: 接收者列表
// subject: 邮件主题
// content: 邮件正文内容
// attachPath: 附件文件路径列表
// reSendCount: 失败重试次数
func (e *EMail) SendEmail(receiver []string, subject string, content string, reSendCount int, attachPath []string) (err error) {
	if reSendCount < 0 {
		return errors.New("reSendCount < 0 is err")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "<"+e.opts.email+">")
	m.SetHeader("To", receiver...)  //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", content) //设置邮件正文

	for _, attach := range attachPath {
		m.Attach(attach) //设置附件, 填写附件路径和名称
	}

	dialer := gomail.NewDialer(e.opts.host, e.opts.port, e.opts.email, e.opts.password)
	for i := 0; i < reSendCount+1; i++ {
		err = dialer.DialAndSend(m)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("send email is err:%s", err.Error())
	}
	return nil
}
