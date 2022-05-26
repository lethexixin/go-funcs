package grpc

import (
	"net"
	"runtime"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	network = "tcp"
)

type Server struct {
	Server *grpc.Server
}

func NewServer() *Server {
	return &Server{
		Server: grpc.NewServer(),
	}
}

type ResisterCallBack func() (s *Server)

func (gs *Server) Service(name string, port string, registerFunc ResisterCallBack) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//起服务
	lis, err := net.Listen(network, ":"+port)
	if err != nil {
		logger.Errorf("failed to listen %s, err:%s", name, err.Error())
		return
	}
	s := registerFunc()
	reflection.Register(s.Server)
	err = s.Server.Serve(lis)
	if err != nil {
		logger.Errorf("failed to start grpc server:%s, err:%s", name, err.Error())
		return
	}
	logger.Infof("grpc server:%s start in: %s", name, port)
}
