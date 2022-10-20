package graceful

import (
	"net/http"
	"sync"
	"time"
)

const defaultShutDownTime = time.Second * 15

var (
	once sync.Once

	healthz     = true
	healthzInfo = HealthzInfo{
		port:   "7299",
		router: "healthz",
	}

	servers = make([]HttpSrvInfo, 0)
)

type HealthzInfo struct {
	port   string
	router string
}

func SetHealthz(port, router string) {
	healthzInfo.port = port
	healthzInfo.router = router
}

type HttpSrvInfo struct {
	*http.Server
	addr string
}

func SetHttpSrvInfo(srv *http.Server, addr string) {
	servers = append(servers, HttpSrvInfo{
		Server: srv,
		addr:   addr,
	})
}
