package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(originData []byte) []byte {
	length := len(originData)
	unpadding := int(originData[length-1])
	return originData[:(length - unpadding)]
}

// AesEncrypt 加密
func AesEncrypt(originData, key []byte) (crypted []byte, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			crypted = nil
			err = fmt.Errorf("%v", exception)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	originData = pkcs7Padding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted = make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return
}

// AesDecrypt 解密
func AesDecrypt(crypted, key []byte) (originData []byte, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			originData = nil
			err = fmt.Errorf("%v", exception)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	originData = pkcs7UnPadding(origData)
	return
}

func demo() {
	key := []byte("1234567890123456") // 16位 key
	result, err := AesEncrypt([]byte("shepard"), key)
	if err != nil {
		panic(err)
	}
	str := base64.StdEncoding.EncodeToString(result)
	fmt.Println(str)

	buf, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	origData, err := AesDecrypt(buf, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(origData))
}
