// Package server implements a querier that limits server in flight queries.
package server

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/loadimpact/resolvent/internal"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type serverLimitingQuerier struct {
	underlying resolvent.Querier
	servers    internal.SemaphoreMap
}

// New returns a querier that limits server in flight queries.
func New(
	underlying resolvent.Querier,
	max uint16,
) (*serverLimitingQuerier, error) {
	if max == 0 {
		return nil, errors.New("invalid max (0): must be positive")
	}
	return &serverLimitingQuerier{
		underlying: underlying,
		servers:    internal.NewSemaphoreMap(max),
	}, nil
}

// Query performs a query when serverwide capacity is available.
func (q *serverLimitingQuerier) Query(
	ctx context.Context,
	protocol resolvent.Protocol,
	local net.IP,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	key, err := constructKey(address, port)
	if err != nil {
		return
	}
	if err = q.servers.Procure(ctx, key); err != nil {
		return
	}
	defer q.servers.Vacate(key)
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
