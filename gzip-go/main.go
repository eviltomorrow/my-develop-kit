package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var path = "/home/shepard/workspace-tmp/local/data.json"
	var result bytes.Buffer
	var writer = gzip.NewWriter(&result)

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var cache bytes.Buffer
	var buf [1024]byte
	var t = 0
	var count = 0
	for {
		n, err := file.Read(buf[t:])
		if err == io.EOF {
			count += n
			cache.Write(buf[:n])
			break
		}
		cache.Write(buf[:n])
		t = 0
		count += n
	}

	writer.Write(cache.Bytes())
	writer.Close()
	fmt.Printf("文件路径： %v, 字节大小： %v \r\n", path, count)
	fmt.Printf("压缩前大小：%v，压缩后大小：%v, 压缩比率： %.2f %%\r\n", count, len(result.Bytes()), 100-float64(len(result.Bytes()))/float64(count)*100)

	reader, err := gzip.NewReader(&result)
	if err != nil {
		log.Fatal(err)
	}

	out, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
