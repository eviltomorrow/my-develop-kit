package main

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
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

	// 写入值
	kv := clientv3.NewKV(client)
	putResp, err := kv.Put(context.Background(), "foo", "Hello world", clientv3.WithPrevKV())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("revision: %v\r\n", putResp.Header.Revision)
	putResp, err = kv.Put(context.Background(), "foo1", "Hello china", clientv3.WithPrevKV())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("revision: %v\r\n", putResp.Header.Revision)

	// 获取
	getResp, err := kv.Get(context.Background(), "foo", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("get: %v\r\n", getResp.Count)
	for _, r := range getResp.Kvs {
		log.Printf("get: %s, %s, %v\r\n", r.Key, r.Value, r.Version)
	}

	// 删除
	putResp, err = kv.Put(context.Background(), "too", "Shepard", clientv3.WithPrevKV())
	if err != nil {
		log.Fatal(err)
	}
	delResp, err := kv.Delete(context.Background(), "too", clientv3.WithFromKey())
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range delResp.PrevKvs {
		log.Printf("delete key is: %s \n Value: %s \n", string(r.Key), string(r.Value))
	}

	// 租期
	lease := clientv3.NewLease(client)
	leaseResp, err := lease.Grant(context.Background(), 10)
	if err != nil {
		log.Fatal(err)
	}
	leaseID := leaseResp.ID

	keepResp, err := lease.KeepAlive(context.Background(), leaseID)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case k := <-keepResp:
				if k == nil {
					log.Printf("租约失效\r\n")
					return
				}
				// log.Printf("lease-id: %v\r\n", k.ID)
			}
		}
	}()
	putResp, err = kv.Put(context.Background(), "student", "kathy", clientv3.WithLease(leaseID))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			kv.Put(context.Background(), "watch", "hi")
			time.Sleep(3 * time.Second)
			kv.Delete(context.Background(), "watch")
			time.Sleep(4 * time.Second)
		}
	}()

	go func() {
		watcher := clientv3.NewWatcher(client)
		watchResp := watcher.Watch(context.Background(), "watch", clientv3.WithPrefix())
		for r := range watchResp {
			for _, event := range r.Events {
				switch event.Type {
				case mvccpb.DELETE:
					log.Printf("Delete Action\r\n")
				case mvccpb.PUT:
					log.Printf("Put Action\r\n")
				}
			}
		}
	}()

	for {
		getResp, err = kv.Get(context.Background(), "student", clientv3.WithPrefix())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("count: %v, value: %v\r\n", getResp.Count, getResp.Kvs)
		if getResp.Count == 0 {
			log.Printf("Over\r\n")
			break
		}

		time.Sleep(2 * time.Second)
	}

}
