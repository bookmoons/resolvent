package resolvent

import (
	"context"
	"net"

	"github.com/miekg/dns"
)

func (r *Resolver) query(
	ctx context.Context,
	address net.IP,
	port uint16,
	qname string,
	qclass uint16,
	qtype uint16,
) (response *dns.Msg, err error) {
	request := new(dns.Msg)
	request.Id = dns.Id()
	request.Question = make([]dns.Question, 1)
	request.Question[0] = dns.Question{
		Name:   dns.Fqdn(qname),
		Qclass: qclass,
		Qtype:  qtype,
	}
	server, err := constructServer(address, port)
	if err != nil {
		return nil, err
	}
	return dns.ExchangeContext(ctx, request, server)
}
