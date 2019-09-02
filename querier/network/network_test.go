package network

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	t.Run("invalid address", func(t *testing.T) {
		querier := New()
		_, _, err := querier.Query(
			context.Background(),
			resolvent.UDP,
			[]byte{1, 2, 3, 4, 5},
			53,
			"epoch.test",
			dns.ClassINET,
			dns.TypeA,
		)
		assert.EqualError(t, err, "invalid IP address")
	})
	t.Run("timeout", func(t *testing.T) {
		querier := New()
		ctx, _ := context.WithTimeout(
			context.Background(),
			100*time.Millisecond,
		)
		_, _, err := querier.Query(
			ctx,
			resolvent.UDP,
			net.ParseIP("192.0.2.1"),
			53,
			"era.test",
			dns.ClassINET,
			dns.TypeA,
		)
		assert.Error(t, err, "incorrect success")
	})
	t.Run("success udp", func(t *testing.T) {
		querier := New()
		answer, err := dns.NewRR("age.test A 192.0.2.2")
		if err != nil {
			t.Fatalf("failed to construct response: %v", err)
		}
		message := &dns.Msg{Answer: []dns.RR{answer}}
		server, port, err := startTestServer(
			t,
			resolvent.UDP,
			[]*dns.Msg{message},
		)
		if err != nil {
			t.Fatalf("failed to start server: %v", err)
		}
		defer server.Shutdown()
		response, _, err := querier.Query(
			context.Background(),
			resolvent.UDP,
			net.ParseIP("127.0.0.1"),
			port,
			"age.test",
			dns.ClassINET,
			dns.TypeA,
		)
		assert.NoError(t, err, "query failed")
		assert.Equal(t, 0, len(response.Ns), "nonempty authority section")
		assert.Equal(t, 0, len(response.Extra), "nonempty additional section")
		assert.Greater(t, len(response.Answer), 0, "empty answer section")
		assert.Less(t, len(response.Answer), 2, "excess answers")
		assert.Equal(
			t,
			"age.test.\t3600\tIN\tA\t192.0.2.2",
			response.Answer[0].String(),
			"incorrect resource record",
		)
	})
	t.Run("success tcp", func(t *testing.T) {
		querier := New()
		answer, err := dns.NewRR("age.test A 192.0.2.3")
		if err != nil {
			t.Fatalf("failed to construct response: %v", err)
		}
		message := &dns.Msg{Answer: []dns.RR{answer}}
		server, port, err := startTestServer(
			t,
			resolvent.TCP,
			[]*dns.Msg{message},
		)
		if err != nil {
			t.Fatalf("failed to start server: %v", err)
		}
		defer server.Shutdown()
		response, _, err := querier.Query(
			context.Background(),
			resolvent.TCP,
			net.ParseIP("127.0.0.1"),
			port,
			"age.test",
			dns.ClassINET,
			dns.TypeA,
		)
		assert.NoError(t, err, "query failed")
		assert.Equal(t, 0, len(response.Ns), "nonempty authority section")
		assert.Equal(t, 0, len(response.Extra), "nonempty additional section")
		assert.Greater(t, len(response.Answer), 0, "empty answer section")
		assert.Less(t, len(response.Answer), 2, "excess answers")
		assert.Equal(
			t,
			"age.test.\t3600\tIN\tA\t192.0.2.3",
			response.Answer[0].String(),
			"incorrect resource record",
		)
	})
}

func startTestServer(
	t *testing.T,
	protocol resolvent.Protocol,
	responses []*dns.Msg,
) (server *dns.Server, port uint16, err error) {
	// Stage responses
	responsesChan := make(chan *dns.Msg, len(responses))
	for _, response := range responses {
		responsesChan <- response
	}
	handler := func(writer dns.ResponseWriter, request *dns.Msg) {
		response := <-responsesChan
		response.Id = request.Id
		err := writer.WriteMsg(response)
		assert.NoError(t, err, "handle error")
	}

	// Start server
	started := make(chan struct{}, 1)
	server = &dns.Server{
		Net:     translateProtocol(protocol),
		Addr:    "127.0.0.1:0",
		Handler: dns.HandlerFunc(handler),
		NotifyStartedFunc: func() {
			started <- struct{}{}
		},
	}
	go func() {
		err := server.ListenAndServe()
		assert.NoError(t, err, "server error")
	}()
	<-started

	// Discover port
	address := extractAddress(server, protocol)
	_, portString, err := net.SplitHostPort(address.String())
	if err != nil {
		return
	}
	portWide, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return
	}
	port = uint16(portWide)
	return
}

func translateProtocol(protocol resolvent.Protocol) (network string) {
	switch protocol {
	case resolvent.UDP:
		return "udp"
	case resolvent.TCP:
		return "tcp"
	default:
		panic("invalid protocol")
	}
}

func extractAddress(
	server *dns.Server,
	protocol resolvent.Protocol,
) (address net.Addr) {
	switch protocol {
	case resolvent.UDP:
		return server.PacketConn.LocalAddr()
	case resolvent.TCP:
		return server.Listener.Addr()
	default:
		panic("invalid protocol")
	}
}
