package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// RsaEncrypt 加密
func RsaEncrypt(data []byte, publicKeyPath string) ([]byte, error) {
	file, err := os.OpenFile(publicKeyPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	key, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	//解密pem格式的公钥
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// RsaDecrypt 解密
func RsaDecrypt(ciphertext []byte, privateKeyPath string) ([]byte, error) {
	file, err := os.OpenFile(privateKeyPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	key, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	//解密
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
func main() {
	var str = "This is shepard, "
	data, err := RsaEncrypt([]byte(str), "rsa_public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("加密前: %s\r\n", str)
	log.Printf("加密后: %s\r\n", base64.StdEncoding.EncodeToString(data))

	origin, err := RsaDecrypt(data, "rsa_private_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("解密后: %s\r\n", origin)
}
