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

type Resolver struct {
	q            querier.Querier
	QueryTimeout time.Duration
}

func New() *Resolver {
	return &Resolver{
		q:            networkQuerier.New(),
		QueryTimeout: defaultQueryTimeout * time.Second,
	}
}

func (r *Resolver) Query(
	ctx context.Context,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, err error) {
	if _, ok := ctx.Deadline(); !ok {
		timed, cancel := context.WithTimeout(ctx, r.QueryTimeout)
		defer cancel()
		ctx = timed
	}
	return r.q.Query(ctx, address, port, qname, qclass, qtype)
}
