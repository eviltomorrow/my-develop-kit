package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/my-develop-kit/grpc-go/pb"
	"google.golang.org/grpc"
)

func TestSimpleGRPC(t *testing.T) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
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

func BenchmarkSimpleGRPC(b *testing.B) {
	conn, _ := grpc.Dial("localhost:8080", grpc.WithInsecure())
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := pb.NewGreeterClient(conn)
		client.SayHello(ctx, &pb.HelloRequest{Name: "shepard"})
	}
	conn.Close()

}
