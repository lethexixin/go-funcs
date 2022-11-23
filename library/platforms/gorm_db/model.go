package gorm_db

import (
	"strings"
)

import (
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormDB struct {
	DB *gorm.DB
}

const (
	SilentLogLevel = "silent"
	ErrorLogLevel  = "error"
	WarnLogLevel   = "warn"
	InfoLogLevel   = "info"
)

type Options struct {
	logLevel        gormLogger.LogLevel
	dsn             string
	maxOpenConn     int
	maxIdleConn     int
	connMaxLifetime int
}

type Option func(*Options)

const (
	DefaultLogLevel        = gormLogger.Warn
	DefaultMysqlDSN        = "root:123456@tcp(localhost:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	DefaultCkDSN           = "http://localhost:8123/test?username=default&password=123456&dial_timeout=200ms&max_execution_time=60"
	DefaultPGDSN           = "host=localhost user=root password=123456 dbname=demo port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	DefaultMaxOpenConn     = 1000
	DefaultMaxIdleConn     = 100
	DefaultConnMaxLifetime = 3600
)

func LogLevel(logLevel string) Option {
	level := DefaultLogLevel
	switch strings.ToLower(logLevel) {
	case SilentLogLevel:
		level = gormLogger.Silent
	case ErrorLogLevel:
		level = gormLogger.Error
	case InfoLogLevel:
		level = gormLogger.Info
	default:
		level = gormLogger.Warn
	}
	return func(o *Options) {
		o.logLevel = level
	}
}

func DSN(dsn string) Option {
	return func(o *Options) {
		o.dsn = dsn
	}
}

func MaxOpenConn(maxOpenConn int) Option {
	return func(o *Options) {
		o.maxOpenConn = maxOpenConn
	}
}

func MaxIdleConn(maxIdleConn int) Option {
	return func(o *Options) {
		o.maxIdleConn = maxIdleConn
	}
}

func ConnMaxLifetime(connMaxLifetime int) Option {
	return func(o *Options) {
		o.connMaxLifetime = connMaxLifetime
	}
}
