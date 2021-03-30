package main

import (
	"context"
	"fmt"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"github.com/skiffer-git/grpc-etcdv3/helloworld"
	"time"
)

func main(){
	ticker := time.NewTicker(2000 * time.Millisecond)
	for t := range ticker.C {
		etcdAddr := "111.52.125.183:2379"  //your etcd svr
		conn := getcdv3.GetConn("sk", etcdAddr,"myrpc")
		if(conn == nil){
			fmt.Println(conn)
			continue
		}
		client := helloworld.NewHelloClient(conn)

		resp, err := client.SayHello(context.Background(), &helloworld.HelloReq{Req: "world"})
		if err == nil {
			fmt.Println( resp.Response,t )

		}else{
			fmt.Println("errrrrrr", err)
		}
	}
}
