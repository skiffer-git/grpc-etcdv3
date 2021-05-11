package main

import (
	"context"
	"fmt"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"github.com/skiffer-git/grpc-etcdv3/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"strconv"
)


type server struct{

}
//SayHello(context.Context, *HelloReq) (*HelloResp, error)
func(s *server) SayHello(ctx context.Context, in *helloworld.HelloReq) (*helloworld.HelloResp, error) {
	pr, _:= peer.FromContext(ctx)
	return &helloworld.HelloResp{Response: in.Req + " from: " + pr.Addr.String()}, nil
}


func work(port int) {
	//
	listener, err := net.Listen("tcp",net.JoinHostPort("127.0.0.1",strconv.Itoa(port) ))
	if err != nil {
		fmt.Println("listen failed")
		return
	}
	//"%s:///%s"
	etcdAddr := "47.112.160.66:2379"
	getcdv3.RegisterEtcd("sk", etcdAddr, "127.0.0.1", port, "myrpc1",10)
	getcdv3.RegisterEtcd4Unique("sk", etcdAddr, "127.0.0.1", port, "myrpc2",10)

	s := grpc.NewServer()
	helloworld.RegisterHelloServer(s, &server{})
	s.Serve(listener)

}

func  main()  {
	//go work(22222)
	work(44444)
}
