package grpclb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// Prefix prefix
var Prefix = "etcd"

var key string
var stopsignal = make(chan struct{}, 1)

// RegisterToETCD register to etcd
func RegisterToETCD(name string, host string, port int, endpoints []string, interval time.Duration, ttl int64) error {
	key = fmt.Sprintf("%s://%s-%s:%d", Prefix, name, host, port)

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
	})

	if err != nil {
		return err
	}

	response, err := client.Grant(context.TODO(), ttl)
	if err != nil {
		return nil
	}

	_, err = client.Put(context.TODO(), key, fmt.Sprintf("%s:%d", host, port), clientv3.WithLease(response.ID))
	if err != nil {
		return err
	}

	keepalive, err := client.KeepAlive(context.TODO(), response.ID)
	if err != nil {
		return err
	}

	go func() {
	loop:
		for {
			select {
			case <-stopsignal:
				break loop
			case _, ok := <-keepalive:
				if !ok {
					break loop
				}
				log.Printf("")
			}
		}

		if _, err := client.Revoke(context.TODO(), response.ID); err != nil {
			return
		}
	}()

	return nil
}

// DestroyFromETCD destroy from etcd
func DestroyFromETCD() error {
	return nil
}
