// Package total implements a querier that limits total in flight queries.
package total

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type totalLimitingQuerier struct {
	underlying resolvent.Querier
	semaphore  chan struct{}
}

// New returns a querier that limits total in flight queries.
func New(
	underlying resolvent.Querier,
	max uint16,
) (*totalLimitingQuerier, error) {
	if max == 0 {
		return nil, errors.New("invalid max (0): must be positive")
	}
	return &totalLimitingQuerier{
		underlying: underlying,
		semaphore:  make(chan struct{}, max),
	}, nil
}

// Query performs a query when resolverwide capacity is available.
func (q *totalLimitingQuerier) Query(
	ctx context.Context,
	protocol resolvent.Protocol,
	local net.IP,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	err = q.procure(ctx)
	if err != nil {
		return
	}
	defer q.vacate()
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

func (q *totalLimitingQuerier) procure(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.semaphore <- struct{}{}:
		return nil
	}
}

func (q *totalLimitingQuerier) vacate() {
	<-q.semaphore
}
