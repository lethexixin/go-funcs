package graceful

import (
	"net/http"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

func unHealthz() {
	healthz = false
}

func handlerUpListen() {
	go func() {
		h := func(w http.ResponseWriter, _ *http.Request) {
			if !healthz {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
		http.HandleFunc("/"+healthzInfo.router, h)

		s := func(w http.ResponseWriter, _ *http.Request) {
			logger.Infof("preStop received")
			// 程序退出, 将 healthz 置为false, k8s的就绪检测才不会将流量分配到该程序
			unHealthz()
		}
		http.HandleFunc("/stop", s)

		if err := http.ListenAndServe(":"+healthzInfo.port, nil); err != nil {
			logger.Errorf("start healthz check err:%s", err.Error())
		}
	}()
}
