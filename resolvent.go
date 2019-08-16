package resolvent

import (
	"context"

	"github.com/miekg/dns"
)

type Resolver struct{}

func (r *Resolver) Query(
	ctx context.Context,
	server string,
	qname string,
	qclass uint16,
	qtype uint16,
) ([]dns.RR, error) {
	message := new(dns.Msg)
	message.Id = dns.Id()
	message.Question = make([]dns.Question, 1)
	message.Question[0] = dns.Question{
		dns.Fqdn(qname),
		qtype,
		qclass,
	}
	response, err := dns.ExchangeContext(ctx, message, server)
	if err != nil {
		return nil, err
	}
	return response.Answer, nil
}
