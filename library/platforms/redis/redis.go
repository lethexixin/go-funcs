package redis

import (
	"context"
	"fmt"
	"time"
)

import (
	"github.com/go-redis/redis/v8"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

type Redis struct {
	Client *redis.Client
}

type Options struct {
	network        string
	host           string
	port           int
	username       string
	password       string
	db             int
	dialTimeoutMs  int
	readTimeoutMs  int
	writeTimeoutMs int
	idleTimeoutMs  int
}

type Option func(*Options)

const (
	DefaultNetwork        = "tcp"
	DefaultHost           = "127.0.0.1"
	DefaultPort           = 6379
	DefaultUsername       = ""
	DefaultPassword       = ""
	DefaultDB             = 0
	DefaultDialTimeoutMs  = 1000 * 5
	DefaultReadTimeoutMs  = 1000 * 3
	DefaultWriteTimeoutMs = 1000 * 3
	DefaultIdleTimeoutMs  = 1000 * 60 * 5
)

func Network(network string) Option {
	return func(o *Options) {
		o.network = network
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

func Username(username string) Option {
	return func(o *Options) {
		o.username = username
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.password = password
	}
}

func DB(db int) Option {
	return func(o *Options) {
		o.db = db
	}
}

func DialTimeoutMs(dialTimeoutMs int) Option {
	return func(o *Options) {
		o.dialTimeoutMs = dialTimeoutMs
	}
}

func ReadTimeoutMs(readTimeoutMs int) Option {
	return func(o *Options) {
		o.readTimeoutMs = readTimeoutMs
	}
}

func WriteTimeoutMs(writeTimeoutMs int) Option {
	return func(o *Options) {
		o.writeTimeoutMs = writeTimeoutMs
	}
}

func IdleTimeoutMs(idleTimeoutMs int) Option {
	return func(o *Options) {
		o.idleTimeoutMs = idleTimeoutMs
	}
}

func (r *Redis) Init(options ...Option) (err error) {
	opts := Options{
		network:        DefaultNetwork,
		host:           DefaultHost,
		port:           DefaultPort,
		username:       DefaultUsername,
		password:       DefaultPassword,
		db:             DefaultDB,
		dialTimeoutMs:  DefaultDialTimeoutMs,
		readTimeoutMs:  DefaultReadTimeoutMs,
		writeTimeoutMs: DefaultWriteTimeoutMs,
		idleTimeoutMs:  DefaultIdleTimeoutMs,
	}

	for _, o := range options {
		o(&opts)
	}

	address := fmt.Sprintf("%s:%d", opts.host, opts.port)
	logger.Infof("create a new redis client, address:%s", address)

	r.Client = redis.NewClient(&redis.Options{
		Network:      opts.network,
		Addr:         address,
		Username:     opts.username,
		Password:     opts.password,
		DB:           opts.db,
		DialTimeout:  time.Duration(opts.dialTimeoutMs) * time.Millisecond,
		ReadTimeout:  time.Duration(opts.readTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(opts.writeTimeoutMs) * time.Millisecond,
		IdleTimeout:  time.Duration(opts.idleTimeoutMs) * time.Millisecond,
	})

	_, err = r.Client.Ping(context.Background()).Result()
	if err != nil {
		logger.Errorf("ping redis client err:%s", err.Error())
		return err
	}

	return nil
}
