package GoDockerTutorialAdapter

import (
	"context"
	"fmt"
	pb "github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	c          pb.GreeterClient
	etcdClient *clientv3.Client
	appId = "com.godockertutorial.tutorial.sayhello"
)

func init() {
	etcdClient, _ = clientv3.New(clientv3.Config{
		Endpoints:   []string{"host.docker.internal:2379"},
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

	resp, err := etcdClient.Get(ctx, appId, clientv3.WithPrefix())
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
