package query

import (
	"context"
	"net"

	"github.com/miekg/dns"
)

type Querier interface {
	Query(
		ctx context.Context,
		address net.IP,
		port uint16,
		qname string,
		qclass uint16,
		qtype uint16,
	) (response *dns.Msg, err error)
}
