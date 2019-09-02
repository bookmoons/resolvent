package double

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
)

type stallQuerier struct {
	Received <-chan struct{}
	received chan struct{}
}

// NewStallQuerier returns a querier taht
func NewStallQuerier() *stallQuerier {
	received := make(chan struct{})
	return &stallQuerier{
		Received: received,
		received: received,
	}
}

// Query never completes. It stalls until canceled or timed out.
// A signal is sent to Received for each query to enable awaiting initiation.
func (q *stallQuerier) Query(
	ctx context.Context,
	protocol resolvent.Protocol,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	q.received <- struct{}{}
	select {
	case <-make(chan struct{}):
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
	panic("escaped stall")
}
