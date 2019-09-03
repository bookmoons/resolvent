package double

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
)

type stubQuerier struct {
	responses []*dns.Msg
}

// NewStubQuerier returns a querier that replies with the specified responses.
func NewStubQuerier(responses []*dns.Msg) *stubQuerier {
	return &stubQuerier{
		responses: responses,
	}
}

// Query returns the next stub response.
func (q *stubQuerier) Query(
	ctx context.Context,
	protocol resolvent.Protocol,
	local net.IP,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	if len(q.responses) == 0 {
		panic("stub responses exhausted")
	}
	response = q.responses[0]
	if len(q.responses) == 1 {
		q.responses = []*dns.Msg{}
	} else {
		q.responses = q.responses[1:]
	}
	return
}
