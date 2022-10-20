package graceful

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"time"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

// UpShutdown 优雅上下线
func UpShutdown() {
	once.Do(func() {
		handlerUpListen()
		handlerShutdown()
	})
}

func handlerShutdown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)

	select {
	case sig := <-signals:
		// 程序退出, 将 healthz 置为false, k8s的就绪检测才不会将流量分配到该程序
		unHealthz()

		logger.Infof("get signal %s, application will shutdown.", sig)
		time.AfterFunc(defaultShutDownTime, func() {
			logger.Warn("shutdown gracefully timeout, application will shutdown immediately. ")
			os.Exit(0)
		})
		beforeShutdown()
		// those signals' original behavior is exit with dump ths stack, so we try to keep the behavior
		for _, dumpSignal := range DumpHeapShutdownSignals {
			if sig == dumpSignal {
				debug.WriteHeapDump(os.Stdout.Fd())
			}
		}
		os.Exit(0)
	}
}

// beforeShutdown provides processing flow before shutdown
func beforeShutdown() {
	destroyAllRequests()
}

func destroyAllRequests() {
	for _, srv := range servers {
		destroyRequest(srv)
	}
}

func destroyRequest(srv HttpSrvInfo) {
	logger.Infof("graceful shutdown --- destroy http server, addr:%s .", srv.addr)
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutDownTime)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("http server shutdown err:%s", err.Error())
	}
}
