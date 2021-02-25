package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"
)

const (
	nmHost  = "host"
	nmPort  = "port"
	nmCerts = "certs"
)

var (
	host  = flag.String("host", "127.0.0.1", "host address")
	port  = flag.String("port", "9090", "server port")
	certs = flag.String("certs", "/var/run/z-console-client/certs", "certs dir")
)

func main() {
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Error: Get current dir failure, nest error: %v\r\n", err)
		return
	}
	if *certs != "" {
		dir = *certs
	} else {
		dir = filepath.Join(dir, "certs")
	}

	cert, err := tls.LoadX509KeyPair(filepath.Join(dir, "client.crt"), filepath.Join(dir, "client.pem"))
	if err != nil {
		log.Printf("Error: LoadX509KeyPair failure, nest error: %v\r\n", err)
		return
	}
	certBytes, err := ioutil.ReadFile(filepath.Join(dir, "ca.crt"))
	if err != nil {
		log.Printf("Error: Read ca cert failure, nest error: %v\r\n", err)
		return
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		log.Printf("Panic: Append ca cert failure\r\n")
		return
	}
	conf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	var dialer = &net.Dialer{
		Timeout: 10 * time.Second,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(*host, *port), conf)
	if err != nil {
		log.Printf("Error: %v\r\n", err)
		return
	}
	defer conn.Close()

	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("Error: Resizing pty: %s\r\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	go func() { _, _ = io.Copy(conn, os.Stdin) }()
	log.Printf("----------------------------------------\r\n")
	log.Printf("- Login server[ %s ] \r\n", *host)
	log.Printf("----------------------------------------\r\n")
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		log.Printf("Error: %v\r\n", err)
	} else {
		log.Printf("----------------------------------------\r\n")
		log.Printf("- Logout server[ %s ] \r\n", *host)
		log.Printf("----------------------------------------\r\n")
	}

}
