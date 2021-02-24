package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/my-develop-kit/grpc-go-etcd/grpclb"
	"github.com/eviltomorrow/my-develop-kit/grpc-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

func main() {
	r := grpclb.NewResolver("test", nil)
	resolver.Register(r)
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", r.Scheme(), "a"), grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Println("connetion ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := pb.NewGreeterClient(conn)
	repley, err := client.SayHello(ctx, &pb.HelloRequest{Name: "shepard"})
	if err != nil {
		log.Fatalf("SayHello error: %v", err)
	}
	fmt.Println(repley.Message)

}
