package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"transfer-file/pb"

	"google.golang.org/grpc"
)

const (
	host = "127.0.0.1"
	port = 8080
)

func main() {
	//建立链接
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTransferServiceClient(conn)

	// 1秒的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.Open(ctx, &pb.FileInfo{
		Path: "/tmp/transfer/test.txt",
		Size: 50,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", resp.Size)
}
