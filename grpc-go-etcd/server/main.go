package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eviltomorrow/my-develop-kit/grpc-go-etcd/pb"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// HelloService hello
type HelloService struct {
}

// SayHello hello
func (hs *HelloService) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello: %d, %s", *port, request.Name)}, nil
}

var port = flag.Int("port", 8080, "Server port")

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpc.EnableTracing = true

	s := grpc.NewServer()
	defer s.Stop()
	defer s.GracefulStop()

	reflection.Register(s)
	pb.RegisterGreeterServer(s, &HelloService{})

	localIP, err := localIP()
	if err != nil {
		log.Fatalf("Get local ip failure, nest error: %v\r\n", err)
	}
	close, err := register("grpclb", localIP, *port, 10)
	if err != nil {
		log.Fatalf("Register grpclb failure, nest error: %v\r\n", err)
	}
	defer close()

	log.Printf("Server start, port: %d\r\n", *port)

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	blockingUntilTermination()
}

func blockingUntilTermination() {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-ch {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	case syscall.SIGUSR1:
	case syscall.SIGUSR2:
	default:
	}
	log.Println("Termination main programming, cleanup function is executed complete")
}

func register(service string, host string, port int, ttl int64) (func(), error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Status(ctx, "localhost:2379")
	if err != nil {
		return nil, err
	}

	leaseResp, err := client.Grant(context.Background(), ttl)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("/%s/%s:%d", service, host, port)
	value := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Register information: key: %s, value: %s", key, value)
	_, err = client.Put(context.Background(), key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}

	keepAlive, err := client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-client.Ctx().Done():
				log.Printf("Etcd server shutdown")
				return
			case k, ok := <-keepAlive:
				if !ok {
					return
				}
				if k == nil {
					log.Printf("Keep alive failure")
				} else {
					log.Printf("Keep alive...")
				}
			}
		}
	}()
	close := func() {
		_, _ = client.Revoke(ctx, leaseResp.ID)
	}

	return close, nil
}

func localIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("unable to determine local ip")
}
