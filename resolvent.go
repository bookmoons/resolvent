// Package resolvent implements DNS resolution.
package resolvent

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent/querier"
	networkQuerier "github.com/loadimpact/resolvent/querier/network"
	"github.com/miekg/dns"
)

const (
	defaultQueryTimeout = 5 // seconds
)

// Resolver implements a DNS resolver.
type Resolver struct {
	q            querier.Querier
	QueryTimeout time.Duration
}

// New returns a DNS resolver.
func New() *Resolver {
	return &Resolver{
		q:            networkQuerier.New(),
		QueryTimeout: defaultQueryTimeout * time.Second,
	}
}

// Query performs an atomic query with a single nameserver.
func (r *Resolver) Query(
	ctx context.Context,
	protocol querier.Protocol,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	if _, ok := ctx.Deadline(); !ok {
		timed, cancel := context.WithTimeout(ctx, r.QueryTimeout)
		defer cancel()
		ctx = timed
	}
	return r.q.Query(ctx, protocol, address, port, qname, qclass, qtype)
}
