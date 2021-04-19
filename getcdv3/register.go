package getcdv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"net"
	"strconv"
	"strings"
)

type RegEtcd struct{
	cli *clientv3.Client
	ctx	context.Context
	cancel context.CancelFunc
	key string
}
var rEtcd *RegEtcd

//    "%s:///%s"
func GetPrefix(schema, serviceName string)(string){
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

//etcdAddr separated by commas
func RegisterEtcd(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int)(error){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","),
	})
	if err != nil {
		//		return fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
		return fmt.Errorf("create etcd clientv3 client failed, errmsg:%v, etcd addr:%s", err, etcdAddr)
	}

	//lease
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := cli.Grant(ctx, int64(ttl))
	if(err != nil){
		return fmt.Errorf("grant failed")
	}

	//  schema:///serviceName/ip:port ->ip:port
	serviceValue := net.JoinHostPort(myHost,  strconv.Itoa(myPort))
	serviceKey := GetPrefix(schema, serviceName)+serviceValue

	//set key->value
	if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("put failed, errmsg:%v， key:%s, value:%s", err, serviceKey, serviceValue)
	}

	//keepalive
	kresp, err := cli.KeepAlive(ctx, resp.ID);
	if  err != nil {
		return fmt.Errorf("keepalive faild, errmsg:%v, lease id:%d", err, resp.ID)
	}

	go func() {
	FLOOP:
		for {
			select {
			case _, ok := <-kresp:
				if ok == true{
				//	fmt.Println("keepalive resp: ", r)
				} else {
					break FLOOP
				}
			}
		}
	}()

	rEtcd = &RegEtcd{ctx: ctx,
		cli:cli,
		cancel:cancel,
		key:serviceKey}

	return nil
}

func UnRegisterEtcd(){
	//delete
	rEtcd.cancel()
	rEtcd.cli.Delete(rEtcd.ctx, rEtcd.key)
}
