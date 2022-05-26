package gorm_db

import (
	"time"
)

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

func (g *GormDB) InitCK(options ...Option) error {
	opts := Options{
		logLevel:        DefaultLogLevel,
		dsn:             DefaultCkDSN,
		maxOpenConn:     DefaultMaxOpenConn,
		maxIdleConn:     DefaultMaxIdleConn,
		connMaxLifetime: DefaultConnMaxLifetime,
	}

	for _, o := range options {
		o(&opts)
	}

	logger.Infof("create a new gorm db - clickhouse, dsn:%s", opts.dsn)

	sql, err := gorm.Open(clickhouse.Open(opts.dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormLogger.Default.LogMode(opts.logLevel),
	})
	if err != nil {
		logger.Errorf("gorm.Open dsn:%s err:%s", opts.dsn, err.Error())
		return nil
	}

	sqlDB, err := sql.DB()
	if err != nil {
		logger.Errorf("gorm.Open dsn:%s err:%s", opts.dsn, err.Error())
		return err
	}

	sqlDB.SetMaxIdleConns(opts.maxIdleConn)
	sqlDB.SetMaxOpenConns(opts.maxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(opts.connMaxLifetime) * time.Second)
	g.DB = sql
	return nil
}
