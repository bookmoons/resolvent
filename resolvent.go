// Package resolvent defines domain types.
package resolvent

import (
	"context"
	"net"
	"time"

	"github.com/miekg/dns"
)

// Protocol is a network transport protocol.
type Protocol int

const (
	UDP Protocol = iota
	TCP
)

// Querier is the interface implemented by DNS queriers.
type Querier interface {
	Query(
		ctx context.Context,
		protocol Protocol,
		local net.IP,
		address net.IP,
		port uint16,
		qname string,
		qclass uint16,
		qtype uint16,
	) (response *dns.Msg, duration time.Duration, err error)
}
