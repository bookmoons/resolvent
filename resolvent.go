package resolvent

import (
	"context"
	"net"

	"github.com/loadimpact/resolvent/query"
	"github.com/loadimpact/resolvent/query/live"
	"github.com/miekg/dns"
)

type Resolver struct {
	querier query.Querier
}

func New() *Resolver {
	return &Resolver{
		querier: &live.LiveQuerier{},
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
	return r.querier.Query(ctx, address, port, qname, qclass, qtype)
}
