package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	nmHTTPPort  = "http-port"
	nmHTTPSPort = "https-port"
	nmCertPath  = "cert-path"
)

var (
	httpPort  = flag.Int(nmHTTPPort, 9090, "HTTP port for service, default: 9090")
	httpsPort = flag.Int(nmHTTPSPort, 9091, "HTTPS port for service, default: 9091")
	certPath  = flag.String(nmCertPath, "/home/shepard/workspace/space-go/project/src/agent/internal/collect/plugins/http/example-server/certs", "Cert path for service, default: ''")
)

func main() {
	flag.Parse()

	var (
		certFile = filepath.Join(*certPath, "server.crt")
		keyFile  = filepath.Join(*certPath, "server.key")
	)

	go func() {
		startupServerHTTP(*httpPort)
	}()
	go func() {
		startupServerHTTPS(*httpsPort, certFile, keyFile)
	}()

	blockingUntilTermination()
}

func blockingUntilTermination() {
	var ch = make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-ch {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	case syscall.SIGUSR1:
	case syscall.SIGUSR2:
	default:
	}
}
