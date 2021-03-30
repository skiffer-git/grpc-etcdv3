package getcdv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
//	"go.etcd.io/etcd/mvcc/mvccpb"

    "github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
	"time"
)

type Resolver struct {
	etcdAddr string
	addrDict map[string]resolver.Address
	cli     *clientv3.Client
	cc      resolver.ClientConn
	serviceName string
	schema string
	watchStartRevision int64
}

var (
	mu      sync.Mutex
	allPrefix map[string]int = make (map[string]int)
)

func (r1 *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

func (r1 *Resolver) Close() {
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error){
	//etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(r.etcdAddr, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("etcd clientv3 client failed: %v, etcd:%s", err, target)
	}

	r.cli  = cli
	r.cc = cc
	r.addrDict =  make(map[string]resolver.Address)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix(r.schema, r.serviceName)
	// get key first
	resp, err := r.cli.Get(ctx, prefix, clientv3.WithPrefix())

	if err == nil {
		for i := range resp.Kvs {
		//	fmt.Println(resp.Kvs[i].Value)
			r.addrDict[string(resp.Kvs[i].Value)] = resolver.Address{Addr: string(resp.Kvs[i].Value)}
		}
		r.update(r.addrDict)
		r.watchStartRevision = resp.Header.Revision + 1
	//	fmt.Println(resp.Header.Revision)
	}else {
		return nil, fmt.Errorf("etcd get failed, prefix: %s", prefix)
	}


	//goroutine watch
	go r.watch(prefix)
	return r, nil
}

func (r *Resolver) update(addrDict map[string]resolver.Address){
	addrList := make([]resolver.Address, 0, len(addrDict))
	for _, v := range addrDict {
		addrList = append(addrList, v)
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func(r*Resolver) Scheme() string{
	return r.schema
}


func (r *Resolver) watch(prefix string) {
	//only one goroutine for same prefix
	mu.Lock()
	_, ok := allPrefix[prefix]
	if ok {
		mu.Unlock()
		return
	} else {
		allPrefix[prefix] = 1
	}
	mu.Unlock()

	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithRev(r.watchStartRevision))
	for n := range rch {
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				(r.addrDict)[string(ev.Kv.Key)] = resolver.Address{Addr: string(ev.Kv.Value)}
			case mvccpb.DELETE:
				delete(r.addrDict, string(ev.Kv.Key))
			}
		}
		r.update(r.addrDict)
	}
}

func GetBuild(schema, etcdaddr, servicename string) (*Resolver){
	r := new(Resolver)
	r.etcdAddr = etcdaddr
	r.schema = schema
	r.serviceName = servicename
	return r
}

func GetConn(schema, etcdaddr, servicename string) (*grpc.ClientConn){
	resolver.Register(GetBuild(schema, etcdaddr, servicename))
	conn, err := grpc.Dial(
		GetPrefix(schema, servicename),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(5)*time.Second),
	)
	if(err != nil) {
		return nil
	}
	return conn
}

