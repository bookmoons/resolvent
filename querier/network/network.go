// Package network implements a querier that performs network exchange.
package network

import (
	"context"
	"net"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/loadimpact/resolvent/internal"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type networkQuerier struct {
	clients map[string]map[resolvent.Protocol]*dns.Client
}

// New returns a querier that performs network exchange.
func New() (querier *networkQuerier, err error) {
	clients, err := constructClients()
	if err != nil {
		return
	}
	querier = &networkQuerier{
		clients: clients,
	}
	return
}

// Query executes an exchange with a single DNS nameserver.
func (q *networkQuerier) Query(
	ctx context.Context,
	protocol resolvent.Protocol,
	local net.IP,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, duration time.Duration, err error) {
	client, err := q.acquireClient(protocol, local)
	if err != nil {
		return
	}
	hostport, err := internal.ConstructHostport(address, port)
	if err != nil {
		return
	}
	request := new(dns.Msg)
	request.Id = dns.Id()
	request.Question = make([]dns.Question, 1)
	request.Question[0] = dns.Question{
		Name:   dns.Fqdn(qname),
		Qclass: qclass,
		Qtype:  qtype,
	}
	return client.ExchangeContext(ctx, request, hostport)
}

func (q *networkQuerier) acquireClient(
	protocol resolvent.Protocol,
	local net.IP,
) (client *dns.Client, err error) {
	addressClients, ok := q.clients[local.String()]
	if !ok {
		err = errors.New("invalid local address")
		return
	}
	client, ok = addressClients[protocol]
	if !ok {
		err = errors.New("invalid protocol")
		return
	}
	return
}
