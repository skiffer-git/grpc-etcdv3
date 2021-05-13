package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"github.com/skiffer-git/grpc-etcdv3/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"strconv"
)

type server struct {
}

//SayHello(context.Context, *HelloReq) (*HelloResp, error)
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloReq) (*helloworld.HelloResp, error) {
	pr, _ := peer.FromContext(ctx)
	return &helloworld.HelloResp{Response: in.Req + " from: " + pr.Addr.String() + "PORT:" + strconv.Itoa(*Port)}, nil
}

func work(port int) {
	//
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *Port))
	if err != nil {
		fmt.Println("listen failed")
		return
	}
	//"%s:///%s"
	etcdAddr := "47.112.160.66:2379"

	getcdv3.RegisterEtcd("sk", etcdAddr, "127.0.0.1", port, "myrpc1", 10)
	//	getcdv3.RegisterEtcd4Unique("sk", etcdAddr, "127.0.0.1", port, "myrpc2", 10)

	s := grpc.NewServer()
	helloworld.RegisterHelloServer(s, &server{})
	s.Serve(listener)

}

var Port = flag.Int("Port", 3000, "listening port")

func main() {
	flag.Parse()
	work(*Port)
}
