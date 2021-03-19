package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/eviltomorrow/my-develop-kit/grpc-go-tls/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	var caFile = "/home/shepard/workspace/space-go/project/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/ca.crt"
	var certFile = "/home/shepard/workspace/space-go/project/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/server.crt"
	var keyFile = "/home/shepard/workspace/space-go/project/src/github.com/eviltomorrow/my-develop-kit/grpc-go-tls/certs/server.pem"
	// Create tls based credential.
	// creds, err := credentials.NewServerTLSFromFile(testdata.Path(certFile), testdata.Path(keyFile))
	// if err != nil {
	// 	log.Fatalf("failed to create credentials: %v", err)
	// }

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("LoadX509KeyPair failure, nest error: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Read ca pem file failure, nest error: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("AppendCertsFromPEM failure")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	})

	s := grpc.NewServer(grpc.Creds(creds))

	// Register EchoServer on the server.
	pb.RegisterGreeterServer(s, &HelloService{})

	log.Println("Server start")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
