package grpclb

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

// Prefix prefix
var Prefix = "etcd"

// Resolver resolver
type Resolver struct {
	endpoints []string
	service   string
}

// NewResolver new resolver
func NewResolver(service string, endpoints []string) resolver.Builder {
	return &Resolver{endpoints: endpoints, service: service}
}

// Build build build
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return nil, nil
}

// Scheme returns the scheme supported by this resolver.
func (r *Resolver) Scheme() string {
	return fmt.Sprintf("%s://%s", Prefix, r.service)
}
