package grpclb

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var stopsignal = make(chan struct{}, 1)

// Register register service
func Register(service string, host string, port int, interval time.Duration, ttl int64, endpoints []string) error {
	val := fmt.Sprintf("%s:%d", host, port)
	key := fmt.Sprintf("%s://%s/%s", Prefix, service, val)

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	response, err := client.Grant(context.TODO(), ttl)
	if err != nil {
		return err
	}

	_, err = client.Put(context.TODO(), key, val, clientv3.WithLease(response.ID))
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
					continue loop
				}
			}
		}

		if _, err := client.Revoke(context.TODO(), response.ID); err != nil {

		}
	}()

	return nil
}

// UnRegister un register
func UnRegister() error {
	stopsignal <- struct{}{}
	stopsignal = make(chan struct{})
	return nil
}
