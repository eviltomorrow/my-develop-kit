package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/trace"
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
	startTrace()

	s := grpc.NewServer(grpc.NumStreamWorkers(3))
	pb.RegisterGreeterServer(s, &HelloService{})

	log.Println("Server start")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	s.Stop()
}

func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}

	go http.ListenAndServe(":8090", nil)
	log.Printf("Trace listen on 8090\r\n")
}
