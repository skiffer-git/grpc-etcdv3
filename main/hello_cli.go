package main

import (
	"context"
	"fmt"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"github.com/skiffer-git/grpc-etcdv3/helloworld"
	"time"
)

func main(){
	ticker := time.NewTicker(5000 * time.Millisecond)
	for t := range ticker.C {

		fmt.Println("start............")
		etcdAddr := "47.112.160.66:2379"  //your etcd svr



		conn1 := getcdv3.GetConn("sk", etcdAddr,"myrpc1")
		//fmt.Println("conn:", conn1)

		client := helloworld.NewHelloClient(conn1)

		resp1, err := client.SayHello(context.Background(), &helloworld.HelloReq{Req: "world"})
		if err == nil {
			fmt.Println("say1:", resp1.Response, t )

		}else{
			fmt.Println("errrrrrr", err)
		}






		conns := getcdv3.GetConn4Unique("sk", etcdAddr, "myrpc2")

		for _, v := range conns {
			conn := v

			client := helloworld.NewHelloClient(conn)

			resp, err := client.SayHello(context.Background(), &helloworld.HelloReq{Req: "world"})
			if err == nil {
				fmt.Println("say2: ", resp.Response,t )

			}else{
				fmt.Println("errrrrrr", err)
			}

		}






	}
}
