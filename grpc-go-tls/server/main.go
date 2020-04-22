package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/eviltomorrow/grpc-go-tls/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
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

	var certFile = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/grpc-go-tls/certs/server.crt"
	var keyFile = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/grpc-go-tls/certs/server.pem"
	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile(testdata.Path(certFile), testdata.Path(keyFile))
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))

	// Register EchoServer on the server.
	pb.RegisterGreeterServer(s, &HelloService{})

	log.Println("Server start")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
