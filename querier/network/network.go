// Package network implements a querier that performs network exchange.
package network

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent/querier"
	"github.com/miekg/dns"
)

type networkQuerier struct {
	udpClient *dns.Client
	tcpClient *dns.Client
}

// New returns a querier that performs network exchange.
func New() *networkQuerier {
	return &networkQuerier{
		udpClient: &dns.Client{
			Net: "udp",
		},
		tcpClient: &dns.Client{
			Net: "tcp",
		},
	}
}

// Query executes an exchange with a single DNS nameserver.
func (q *networkQuerier) Query(
	ctx context.Context,
	protocol querier.Protocol,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	hostport, err := constructHostport(address, port)
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
	if protocol == querier.TCP {
		return q.tcpClient.ExchangeContext(ctx, request, hostport)
	}
	return q.udpClient.ExchangeContext(ctx, request, hostport)
}
