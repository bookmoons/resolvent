// Package route implements a querier that limits route in flight queries.
package route

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/loadimpact/resolvent/internal"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type routeLimitingQuerier struct {
	underlying resolvent.Querier
	semaphore  internal.SemaphoreMap
}

// New returns a querier that limits route in flight queries.
func New(
	underlying resolvent.Querier,
	max uint16,
) (*routeLimitingQuerier, error) {
	if max == 0 {
		return nil, errors.New("invalid max (0): must be positive")
	}
	return &routeLimitingQuerier{
		underlying: underlying,
		semaphore:  internal.NewSemaphoreMap(max),
	}, nil
}

// Query performs a query when routewide capacity is available.
func (q *routeLimitingQuerier) Query(
	ctx context.Context,
	protocol resolvent.Protocol,
	local net.IP,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	key, err := constructKey(local, address)
	if err != nil {
		return
	}
	if err = q.semaphore.Procure(ctx, key); err != nil {
		return
	}
	defer q.semaphore.Vacate(key)
	return q.underlying.Query(
		ctx,
		protocol,
		local,
		address,
		port,
		qname,
		qclass,
		qtype,
	)
}
