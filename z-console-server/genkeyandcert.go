package server

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

// ApplicationInformation 申请信息
type ApplicationInformation struct {
	CertificateConfig    *CertificateConfig
	CommonName           string
	CountryName          string
	ProvinceName         string
	LocalityName         string
	OrganizationName     string
	OrganizationUnitName string
}

// CertificateConfig 证书配置
type CertificateConfig struct {
	IsCA           bool
	IP             []net.IP
	DNS            []string
	ExpirationTime time.Duration
}

// GenerateCertificate 生成证书(私钥，证书)
func GenerateCertificate(caKey *rsa.PrivateKey, caCert *x509.Certificate, bits int, info *ApplicationInformation) ([]byte, []byte, error) {
	if !info.CertificateConfig.IsCA {
		if caKey == nil || caCert == nil {
			return nil, nil, fmt.Errorf("Miss ca key/cert")
		}
	}

	priv, err := rsa.GenerateKey(cryptorand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	var template = x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         fmt.Sprintf("%s", info.CommonName),
			Country:            []string{info.CountryName},
			Province:           []string{info.ProvinceName},
			Locality:           []string{info.LocalityName},
			Organization:       []string{info.OrganizationName},
			OrganizationalUnit: []string{info.OrganizationUnitName},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(info.CertificateConfig.ExpirationTime),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
	}

	if info.CertificateConfig.IsCA {
		template.IsCA = true
	} else {
		if i := net.ParseIP(info.CommonName); i != nil {
			template.IPAddresses = append(template.IPAddresses, i)
		} else {
			template.DNSNames = append(template.DNSNames, info.CommonName)
		}
		template.IPAddresses = append(template.IPAddresses, info.CertificateConfig.IP...)
		template.DNSNames = append(template.DNSNames, info.CertificateConfig.DNS...)
	}

	var key *rsa.PrivateKey

	if info.CertificateConfig.IsCA {
		caCert = &template
		key = priv
	} else {
		key = caKey
	}

	certBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, caCert, &priv.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	return x509.MarshalPKCS1PrivateKey(priv), certBytes, nil
}

// ReadCertificate 读取证书
func ReadCertificate(path string) (*x509.Certificate, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, fmt.Errorf("Decode certificate failure, block is nil")
	}

	return x509.ParseCertificate(block.Bytes)
}

// WriteCertificate 写出证书
func WriteCertificate(path string, cert []byte) error {
	_, err := x509.ParseCertificate(cert)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := pem.Encode(&buffer, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	return err
}

// ReadPKCS1PrivateKey 读取 PKCS1 私钥
func ReadPKCS1PrivateKey(path string) (*rsa.PrivateKey, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, fmt.Errorf("Decode private key failure, block is nil")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ReadPKCS8PrivateKey 读取 PKCS8 私钥
func ReadPKCS8PrivateKey(path string) (*rsa.PrivateKey, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, fmt.Errorf("Decode private key failure, block is nil")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if val, ok := key.(*rsa.PrivateKey); ok {
		return val, nil
	}
	return nil, fmt.Errorf("ParsePKCS8PrivateKey failure")
}

// WritePKCS1PrivateKey 写出 PKCS! 私钥
func WritePKCS1PrivateKey(path string, privKey []byte) error {
	_, err := x509.ParsePKCS1PrivateKey(privKey)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := pem.Encode(&buffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKey}); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	return err
}

// WritePKCS8PrivateKey 写出 PKCS8 私钥
func WritePKCS8PrivateKey(path string, privKey []byte) error {
	priv, err := x509.ParsePKCS1PrivateKey(privKey)
	if err != nil {
		return err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := pem.Encode(&buffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	return err
}
