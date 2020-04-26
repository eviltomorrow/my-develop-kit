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

	
	var caCertFile = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/grpc-go-tls/certs/ca.crt"
	var clientCertFile = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/grpc-go-tls/certs/client.crt"
	var clientCertKey = "/home/shepard/workspace-agent/project-go/src/github.com/eviltomorrow/grpc-go-tls/certs/client.pem"
	var serverName = "localhost"

	// Load the certificates from disk
	certificate, err := tls.LoadX509KeyPair(clientCertFile, clientCertKey)
	if err != nil {
		log.Fatalf("LoadX509KeyPair failure: %v", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Load caCertFile failure: %v", err)
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("AppendCertsFromPEM failure: %v", err)
	}

	// Create the TLS credentials for transport
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	// Set up a connection to the server.
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(creds))
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
