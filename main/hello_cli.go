package main

import (
	"context"
	"fmt"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"github.com/skiffer-git/grpc-etcdv3/helloworld"
	"time"
)

func main() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	etcdAddr := "47.112.160.66:2379" //your etcd svr
	//	conn1 := getcdv3.GetConn("sk", etcdAddr, "myrpc1")
	for t := range ticker.C {

		fmt.Println("start............")

		//		fmt.Println("conn:", conn1)

		p, eee := getcdv3.GetConnPool("sk", etcdAddr, "myrpc1")
		if eee != nil {
			continue
		}
		conn1 := p.ClientConn
		if conn1 == nil {
			fmt.Println("get client failed")
		}
		client := helloworld.NewHelloClient(conn1)

		resp1, err := client.SayHello(context.Background(), &helloworld.HelloReq{Req: "world"})
		p.Close()
		if err == nil {
			fmt.Println("say1:", resp1.Response, t)

		} else {
			fmt.Println("errrrrrr", err)
		}

		/*
			time.Sleep(1000 * time.Second)
			conns := getcdv3.GetConn4Unique("sk", etcdAddr, "myrpc2")

			for _, v := range conns {
				conn := v

				client := helloworld.NewHelloClient(conn)

				resp, err := client.SayHello(context.Background(), &helloworld.HelloReq{Req: "world"})
				conn.Close()

				if err == nil {
					fmt.Println("say2: ", resp.Response, t)

				} else {
					fmt.Println("errrrrrr", err)
				}

			}

		*/

	}

}
