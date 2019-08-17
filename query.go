package resolvent

import (
	"context"

	"github.com/miekg/dns"
)

func (r *Resolver) query(
	ctx context.Context,
	server string,
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
	return dns.ExchangeContext(ctx, request, server)
}
