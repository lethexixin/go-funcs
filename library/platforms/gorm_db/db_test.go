package gorm_db

import (
	"os"
	"testing"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

func TestCK(t *testing.T) {
	dsn := "http://localhost:8123/test?username=default&password=123456"
	ck := new(GormDB)
	err := ck.InitCK(DSN(dsn),
		LogLevel("warn"))
	if err != nil {
		logger.Fatalf("failed to create clickhouse connected:%s", err.Error())
		os.Exit(-1)
	}

	ver := ""
	ck.DB.Raw("select version()").Take(&ver)
	t.Log(ver)
}
