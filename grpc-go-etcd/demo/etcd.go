package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Status(ctx, "localhost:2379")
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Get(context.Background(), "foo", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	for _, kv := range response.Kvs {
		fmt.Printf("key: %s, val: %s\r\n", kv.Key, kv.Value)
	}

}
