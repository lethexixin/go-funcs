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
		if err := http.ListenAndServe(":"+healthzInfo.port, nil); err != nil {
			logger.Fatalf("start healthz check err:%s", err.Error())
		}
	}()
}
