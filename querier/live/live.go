package live

import (
	"context"
	"net"

	"github.com/miekg/dns"
)

type liveQuerier struct {
	client *dns.Client
}

func New() *liveQuerier {
	client := new(dns.Client)
	return &liveQuerier{
		client: client,
	}
}

func (q *liveQuerier) Query(
	ctx context.Context,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, err error) {
	server, err := constructServer(address, port)
	if err != nil {
		return nil, err
	}
	request := new(dns.Msg)
	request.Id = dns.Id()
	request.Question = make([]dns.Question, 1)
	request.Question[0] = dns.Question{
		Name:   dns.Fqdn(qname),
		Qclass: qclass,
		Qtype:  qtype,
	}
	response, _, err = q.client.ExchangeContext(ctx, request, server)
	return response, err
}
