package GoDockerTutorialAdapter

import (
	"context"
	pb "github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "172.17.0.2:9090"
)

var (
	c pb.GreeterClient
)

func init() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c = pb.NewGreeterClient(conn)
}

func SayHello(ctx context.Context, req *pb.HelloRequest) *pb.HelloReply {
	rr, err := c.SayHello(ctx, &pb.HelloRequest{Name: "JTF"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return rr
}
