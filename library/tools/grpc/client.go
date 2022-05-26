package grpc

import (
	"runtime"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

import (
	"google.golang.org/grpc"
)

const (
	server = "127.0.0.1"
)

type Client struct {
	ClientConn *grpc.ClientConn
}

func NewGrpcClient(name string, port string) *Client {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//建立连接
	conn, err := grpc.Dial(server+":"+port, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("failed to dial grpc server:%s, err:%s", name, err.Error())
		return nil
	}
	return &Client{
		ClientConn: conn,
	}
}

func (gc Client) Close() {
	if gc.ClientConn != nil {
		_ = gc.ClientConn.Close()
	}
}
