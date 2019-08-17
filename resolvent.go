package resolvent

import (
	"context"

	"github.com/miekg/dns"
)

type Resolver struct{}

func (r *Resolver) Query(
	ctx context.Context,
	server string,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, err error) {
	return r.query(ctx, server, qname, qclass, qtype)
}
