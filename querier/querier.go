// Package querier defines the interface implemented by DNS queriers.
package querier

import (
	"context"
	"net"
	"time"

	"github.com/miekg/dns"
)

// Querier is the interface implemented by DNS queriers.
type Querier interface {
	Query(
		ctx context.Context,
		address net.IP,
		port uint16,
		qname string,
		qclass uint16,
		qtype uint16,
	) (response *dns.Msg, duration time.Duration, err error)
}
