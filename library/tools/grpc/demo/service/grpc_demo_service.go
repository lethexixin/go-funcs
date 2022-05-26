package main

import (
	"net"
	"runtime"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
	demo "github.com/lethexixin/go-funcs/library/tools/grpc/demo"
)

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	network     = "tcp"
	portService = "39090"
)

type GrpcDemoServer struct {
	returnName   []*demo.Name // read-only after initialized
	returnResult []*demo.Result
}

func (s *GrpcDemoServer) GetName(ctx context.Context, in *demo.Person) (*demo.Name, error) {
	return &demo.Name{
		Message: "hello :" + in.FirstName + in.LastName,
	}, nil
}

func (s *GrpcDemoServer) AddOperation(ctx context.Context, in *demo.Param) (*demo.Result, error) {
	return &demo.Result{
		Message: in.X + in.Y,
	}, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// start server
	lis, err := net.Listen(network, ":"+portService)
	if err != nil {
		logger.Fatalf("failed to listen port:%s, err:%s", portService, err.Error())
	}
	s := grpc.NewServer()
	server := &GrpcDemoServer{}
	demo.RegisterGrpcDemoServiceServer(s, server)
	reflection.Register(s)
	logger.Infof("demo grpc server start in port:%s", portService)
	err = s.Serve(lis)
	if err != nil {
		logger.Errorf("failed to start grpc server. err:%s", err.Error())
		return
	}
}
