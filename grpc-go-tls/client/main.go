package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/grpc-go-tls/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

func main() {

	var certFile = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/grpc-go-tls/certs/ca.crt"
	var serverName = "localhost"

	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile(testdata.Path(certFile), serverName)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(creds))
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
