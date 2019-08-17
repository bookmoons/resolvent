package resolvent

import (
	"context"
	"net"

	"github.com/miekg/dns"
)

type Resolver struct{}

func (r *Resolver) Query(
	ctx context.Context,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, err error) {
	return r.query(ctx, address, port, qname, qclass, qtype)
}
