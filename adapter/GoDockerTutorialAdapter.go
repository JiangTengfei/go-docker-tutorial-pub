package GoDockerTutorialAdapter

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	pb "github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	c          pb.GreeterClient
	etcdClient *clientv3.Client
)

func init() {
	etcdClient, _ = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	ip, _ := getServiceIp()

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c = pb.NewGreeterClient(conn)

}

func getServiceIp() (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	resp, err := etcdClient.Get(ctx, "server_ip")
	cancel()
	if err != nil {
		fmt.Errorf("error while put: %+v", err)
		return "", err
	}
	log.Printf("resp: %+v", resp.Kvs)
	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}
	return "", fmt.Errorf("etcd key: %s not found", "server_ip")
}

func SayHello(ctx context.Context, req *pb.HelloRequest) *pb.HelloReply {
	rr, err := c.SayHello(ctx, &pb.HelloRequest{Name: req.Name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return rr
}
