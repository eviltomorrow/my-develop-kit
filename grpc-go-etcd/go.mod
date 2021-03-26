module github.com/eviltomorrow/my-develop-kit/grpc-go-etcd

go 1.15

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/eviltomorrow/my-develop-kit v0.0.0-20210326075441-806b758725c1
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.1
	github.com/prometheus/client_golang v1.10.0 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/grpc v1.36.1
	google.golang.org/grpc/examples v0.0.0-20210326035646-2456c5cff04b // indirect
)
