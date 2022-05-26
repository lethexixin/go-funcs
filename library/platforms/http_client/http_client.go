package http_client

import (
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	Client *http.Client
}

type Options struct {
	dialTimeout        int
	dialKeepAlive      int
	maxIdleConn        int
	idleConnTimeout    int
	maxIdleConnPerHost int
	timeout            int
}

type Option func(*Options)

func DialTimeout(dialTimeout int) Option {
	return func(o *Options) {
		o.dialTimeout = dialTimeout
	}
}

func DialKeepAlive(dialKeepAlive int) Option {
	return func(o *Options) {
		o.dialKeepAlive = dialKeepAlive
	}
}

func MaxIdleConn(maxIdleConn int) Option {
	return func(o *Options) {
		o.maxIdleConn = maxIdleConn
	}
}

func IdleConnTimeout(idleConnTimeout int) Option {
	return func(o *Options) {
		o.idleConnTimeout = idleConnTimeout
	}
}

func MaxIdleConnPerHost(maxIdleConnPerHost int) Option {
	return func(o *Options) {
		o.maxIdleConnPerHost = maxIdleConnPerHost
	}
}

func Timeout(timeout int) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

var (
	DefaultDialTimeout        = 15
	DefaultDialKeepAlive      = 15
	DefaultMaxIdleConn        = 100
	DefaultIdleConnTimeout    = 15
	DefaultMaxIdleConnPerHost = 100
	DefaultTimeout            = 15
)

func (h *HttpClient) Init(options ...Option) {
	opts := Options{
		dialTimeout:        DefaultDialTimeout,
		dialKeepAlive:      DefaultDialKeepAlive,
		maxIdleConn:        DefaultMaxIdleConn,
		idleConnTimeout:    DefaultIdleConnTimeout,
		maxIdleConnPerHost: DefaultMaxIdleConnPerHost,
		timeout:            DefaultTimeout,
	}

	for _, o := range options {
		o(&opts)
	}

	h.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(opts.dialTimeout) * time.Second,
				KeepAlive: time.Duration(opts.dialKeepAlive) * time.Second,
			}).DialContext,
			MaxIdleConns:        opts.maxIdleConn,
			IdleConnTimeout:     time.Duration(opts.idleConnTimeout) * time.Second,
			MaxIdleConnsPerHost: opts.maxIdleConnPerHost,
		},
		Timeout: time.Duration(opts.timeout) * time.Second,
	}
}
