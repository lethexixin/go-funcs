package graceful

import (
	"net/http"
	"sync"
	"time"
)

import (
	"github.com/labstack/echo/v4"
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

	echos = make([]EchoInfo, 0)
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

type EchoInfo struct {
	*echo.Echo
	addr string
}

func SetEchoInfo(echo *echo.Echo, addr string) {
	echos = append(echos, EchoInfo{
		Echo: echo,
		addr: addr,
	})
}
