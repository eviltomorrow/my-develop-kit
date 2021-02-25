package server

import (
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

const (
	nmRebuild = "rebuild"
	nmPort    = "port"
	nmCerts   = "certs"
)

var (
	certs   = flag.String(nmCerts, "", "certs dir")
	rebuild = flag.Bool(nmRebuild, false, "rebuild certs")
	port    = flag.Int(nmPort, 9090, "TCP Server Port")
)

var cfg = DefaultGlobalConfig
var cpf []func() error

// Main main
func Main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Unknown error: %v\r\n", err)
			debug.PrintStack()
			log.Printf("Stack trace: %s\r\n", string(debug.Stack()))
		}
	}()
	flag.Parse()
	cfg.Load(overrideFlags)

	if *rebuild {
		rebuildCertsAndKey()
		os.Exit(0)
	}

	printInfo()
	checkCertsAndKey()
	registerCleanupFunc()
	startupApplication()
	blockingUntilTermination()
}

func printInfo() {
	log.Printf("Config information: %s\r\n", cfg.String())
}

func overrideFlags(cfg *Config) {
	actualFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		actualFlags[f.Name] = true
	})

	if actualFlags[nmPort] {
		cfg.System.Port = *port
	}
	if actualFlags[nmCerts] {
		cfg.System.CertsDir = *certs
	}
}

func rebuildCertsAndKey() {
	if err := os.MkdirAll(cfg.System.CertsDir, 0755); err != nil {
		log.Printf("Fatal: MkdirAll failure, dir: %s, nest error: %v\r\n", cfg.System.CertsDir, err)
		os.Exit(1)
	}
	caPrivBytes, caCertBytes, err := GenerateCertificate(nil, nil, 2048, &ApplicationInformation{
		CertificateConfig: &CertificateConfig{
			IsCA:           true,
			ExpirationTime: 24 * time.Hour * 365 * 3,
		},
		CommonName:           "z-console-server.com",
		CountryName:          "CN",
		ProvinceName:         "BeiJing",
		LocalityName:         "BeiJing",
		OrganizationName:     "Eviltomorrow Inc",
		OrganizationUnitName: "Development",
	})
	if err != nil {
		log.Printf("Error: GenerateCertificate ca cert failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	if err := WritePKCS1PrivateKey(filepath.Join(cfg.System.CertsDir, "ca.key"), caPrivBytes); err != nil {
		log.Printf("Error: WritePKCS1PrivateKey ca key failure, nest error: %v\r\n", err)
		os.Exit(1)
	}

	if err := WriteCertificate(filepath.Join(cfg.System.CertsDir, "ca.crt"), caCertBytes); err != nil {
		log.Printf("Error: WriteCertificate ca cert failure, nest error: %v\r\n", err)
		os.Exit(1)
	}

	// server
	caKey, err := x509.ParsePKCS1PrivateKey(caPrivBytes)
	if err != nil {
		log.Printf("Error: ParsePKCS1PrivateKey CA key failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		log.Printf("Error: ParseCertificate CA cert failure, nest error: %v\r\n", err)
		os.Exit(1)
	}

	getLocalIP := func() (string, error) {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return "", err
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().String()
		idx := strings.LastIndex(localAddr, ":")
		return localAddr[0:idx], nil
	}

	localIP, err := getLocalIP()
	if err != nil {
		log.Printf("Error: GetLocalIP failure, nest error: %v\r\n", err)
		os.Exit(1)
	}

	serverPrivBytes, serverCertBytes, err := GenerateCertificate(caKey, caCert, 2048, &ApplicationInformation{
		CertificateConfig: &CertificateConfig{
			IP:             []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP(localIP)},
			ExpirationTime: 24 * time.Hour * 365,
		},
		CommonName:           "localhost",
		CountryName:          "CN",
		ProvinceName:         "BeiJing",
		LocalityName:         "BeiJing",
		OrganizationName:     "Eviltomorrow Inc",
		OrganizationUnitName: "Development",
	})
	if err != nil {
		log.Printf("GenerateCertificate server cert failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	// WritePKCS1PrivateKey(filepath.Join(cfg.System.CertsDir, "server.key"), serverPrivBytes)
	if err := WritePKCS8PrivateKey(filepath.Join(cfg.System.CertsDir, "server.pem"), serverPrivBytes); err != nil {
		log.Printf("Error: WritePKCS8PrivateKey server pem failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	if err := WriteCertificate(filepath.Join(cfg.System.CertsDir, "server.crt"), serverCertBytes); err != nil {
		log.Printf("Error: WritePKCS1PrivateKey server cert failure, nest error: %v\r\n", err)
		os.Exit(1)
	}

	// client
	clientPrivBytes, clientCertBytes, err := GenerateCertificate(caKey, caCert, 2048, &ApplicationInformation{
		CertificateConfig: &CertificateConfig{
			ExpirationTime: 24 * time.Hour * 365,
		},
		CommonName:           "localhost",
		CountryName:          "CN",
		ProvinceName:         "BeiJing",
		LocalityName:         "BeiJing",
		OrganizationName:     "Apple Inc",
		OrganizationUnitName: "Dev",
	})
	if err != nil {
		log.Printf("Error: GenerateCertificate client key failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	if err := WritePKCS8PrivateKey(filepath.Join(cfg.System.CertsDir, "client.pem"), clientPrivBytes); err != nil {
		log.Printf("Error: WritePKCS8PrivateKey client pem failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
	if err := WriteCertificate(filepath.Join(cfg.System.CertsDir, "client.crt"), clientCertBytes); err != nil {
		log.Printf("Error: WritePKCS8PrivateKey client pem failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
}

func checkCertsAndKey() {
	existfile := func(path string) error {
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("Path is dir")
		}
		return nil
	}

	var caCert = filepath.Join(cfg.System.CertsDir, "ca.crt")
	if err := existfile(caCert); err != nil {
		log.Printf("Error: Ca cert is not exist, certs-dir: %s\r\n", caCert)
		os.Exit(1)
	}

	var serverCert = filepath.Join(cfg.System.CertsDir, "server.crt")
	if err := existfile(serverCert); err != nil {
		log.Printf("Error: Server cert is not exist, certs-dir: %s\r\n", cfg.System.CertsDir)
		os.Exit(1)
	}

	var serverPem = filepath.Join(cfg.System.CertsDir, "server.pem")
	if err := existfile(serverPem); err != nil {
		log.Printf("Server pem is not exist, certs-dir: %s\r\n", cfg.System.CertsDir)
		os.Exit(1)
	}
}

func registerCleanupFunc() {
}

func startupApplication() {
	if err := StartupServer(); err != nil {
		log.Printf("Error: startup tcp server failure, nest error: %v\r\n", err)
		os.Exit(1)
	}
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
	for _, f := range cpf {
		f()
	}
}
