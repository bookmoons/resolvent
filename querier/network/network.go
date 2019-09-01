// Package network implements a querier that performs network exchange.
package network

import (
	"context"
	"net"
	"time"

	"github.com/miekg/dns"
)

type networkQuerier struct {
	client *dns.Client
}

// New returns a querier that performs network exchange.
func New() *networkQuerier {
	client := new(dns.Client)
	return &networkQuerier{
		client: client,
	}
}

// Query executes an exchange with a single DNS nameserver.
func (q *networkQuerier) Query(
	ctx context.Context,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	server, err := constructServer(address, port)
	if err != nil {
		return nil, 0, err
	}
	request := new(dns.Msg)
	request.Id = dns.Id()
	request.Question = make([]dns.Question, 1)
	request.Question[0] = dns.Question{
		Name:   dns.Fqdn(qname),
		Qclass: qclass,
		Qtype:  qtype,
	}
	return q.client.ExchangeContext(ctx, request, server)
}
