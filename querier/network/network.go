package network

import (
	"context"
	"net"

	"github.com/miekg/dns"
)

type networkQuerier struct {
	client *dns.Client
}

func New() *networkQuerier {
	client := new(dns.Client)
	return &networkQuerier{
		client: client,
	}
}

func (q *networkQuerier) Query(
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
