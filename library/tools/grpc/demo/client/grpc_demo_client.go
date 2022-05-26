package main

import (
	"runtime"
	"sync"
	"time"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
	demo "github.com/lethexixin/go-funcs/library/tools/grpc/demo"
)

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	wg sync.WaitGroup
)

const (
	network     = "tcp"
	server      = "127.0.0.1"
	portService = "39090"
	parallel    = 50   //连接并行度
	times       = 1000 //每连接请求次数
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	currTime := time.Now()

	//并行请求
	for i := 0; i < int(parallel); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exe()
		}()
	}
	wg.Wait()

	logger.Infof("time taken: %.2f s", time.Now().Sub(currTime).Seconds())
}

func exe() {
	//建立连接
	conn, err := grpc.Dial(server+":"+portService, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("start grpc connect err:%s", err.Error())
		return
	}
	defer conn.Close()

	client := demo.NewGrpcDemoServiceClient(conn)

	for i := 0; i < int(times); i++ {
		testGetName(client)
	}

	testAdd(client)
}

func testGetName(client demo.GrpcDemoServiceClient) {
	req := &demo.Person{
		FirstName: "xin",
		LastName:  "xi",
	}

	result, _ := client.GetName(context.Background(), req)
	if nil != result {
		logger.Info(result.Message)
	} else {
		logger.Info("response is nil")
	}
}

func testAdd(client demo.GrpcDemoServiceClient) {
	req := &demo.Param{
		X: 5,
		Y: 6,
	}

	result, _ := client.AddOperation(context.Background(), req)
	if nil != result {
		logger.Info(result.Message)
	} else {
		logger.Info("response is nil")
	}
}
