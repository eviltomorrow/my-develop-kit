package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
)

// StartupServer startup server
func StartupServer() error {
	cert, err := tls.LoadX509KeyPair(filepath.Join(cfg.System.CertsDir, "server.crt"), filepath.Join(cfg.System.CertsDir, "server.pem"))
	if err != nil {
		return err
	}

	caCert, err := ioutil.ReadFile(filepath.Join(cfg.System.CertsDir, "ca.crt"))
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return fmt.Errorf("Append ca cert failure")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}
	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", cfg.System.Port), config)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error: Accept connection failure, nest error: %v\r\n", err)
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(conn net.Conn) {
	defer conn.Close()
	buildPTY(conn)
}
