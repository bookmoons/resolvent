// Package resolventtest provides resolvent test helpers.
package resolventtest

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

// CancelableQuery initiates a simple cancelable query.
func CancelableQuery(
	ctx context.Context,
	querier resolvent.Querier,
) (response *dns.Msg, err error) {
	response, _, err = querier.Query(
		ctx,
		resolvent.UDP,
		net.ParseIP("192.0.2.1"),
		53,
		"age.test",
		dns.ClassINET,
		dns.TypeA,
	)
	return
}

// DeepEqual asserts that two values are deeply equal.
func DeepEqual(
	t *testing.T,
	expected interface{},
	actual interface{},
	msgAndArgs ...interface{},
) {
	if len(msgAndArgs) == 0 {
		msgAndArgs = []interface{}{"values not deeply equal"}
	}
	equal := cmp.Equal(expected, actual)
	require.True(t, equal, msgAndArgs...)
}

// MakeMessages constructs a slice of simple DNS messages.
func MakeMessages(
	t *testing.T,
	rrs []string,
	msgAndArgs ...interface{},
) (messages []*dns.Msg) {
	if len(msgAndArgs) == 0 {
		msgAndArgs = []interface{}{"failed to make message"}
	}
	messages = make([]*dns.Msg, len(rrs))
	for i, rr := range rrs {
		answer, err := dns.NewRR(rr)
		require.NoError(t, err, msgAndArgs...)
		messages[i] = &dns.Msg{Answer: []dns.RR{answer}}
	}
	return
}

// SimpleQuery exercises Querier.Query with simple arguments.
func SimpleQuery(
	t *testing.T,
	querier resolvent.Querier,
	msgAndArgs ...interface{},
) (response *dns.Msg) {
	response, _, err := querier.Query(
		context.Background(),
		resolvent.UDP,
		net.ParseIP("192.0.2.1"),
		53,
		"epoch.test",
		dns.ClassINET,
		dns.TypeA,
	)
	require.NoError(t, err, msgAndArgs...)
	return
}

// TimedQuery initiates a simple query with a timeout.
func TimedQuery(
	querier resolvent.Querier,
	timeout time.Duration,
) (response *dns.Msg, err error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	response, _, err = querier.Query(
		ctx,
		resolvent.UDP,
		net.ParseIP("192.0.2.1"),
		53,
		"era.test",
		dns.ClassINET,
		dns.TypeA,
	)
	return
}
