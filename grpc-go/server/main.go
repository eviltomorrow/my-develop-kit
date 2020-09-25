package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/eviltomorrow/my-develop-kit/grpc-go/pb"
)

// HelloService hello
type HelloService struct {
}

// SayHello hello
func (hs *HelloService) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + request.Name}, nil
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpc.EnableTracing = true
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &HelloService{})

	log.Println("Server start")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	s.Stop()
}
