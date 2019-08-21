package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructServer(t *testing.T) {
	t.Parallel()
	pass := []struct {
		address net.IP
		port    uint16
		result  string
	}{
		{net.ParseIP("192.0.2.1"), 49152, "192.0.2.1:49152"},
		{net.ParseIP("2001:db8::1"), 49153, "[2001:db8::1]:49153"},
	}
	for _, item := range pass {
		t.Run(item.address.String(), func(t *testing.T) {
			server, err := constructServer(item.address, item.port)
			assert.NoError(t, err, "incorrect fail")
			assert.Equal(t, server, item.result)
		})
	}
	t.Run("invalid", func(t *testing.T) {
		var address net.IP = []byte{1, 2, 3, 4, 5}
		_, err := constructServer(address, 50000)
		assert.EqualError(t, err, "invalid IP address")
	})
}

func TestIsIPv4(t *testing.T) {
	t.Parallel()
	pass := []net.IP{
		net.ParseIP("192.0.2.1"),
		net.ParseIP("198.51.100.1"),
		net.ParseIP("203.0.113.1"),
		net.ParseIP("::ffff:192.0.2.2"),
	}
	fail := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::88"),
		[]byte{1, 2, 3, 4, 5},
	}
	for _, item := range pass {
		t.Run(item.String(), func(t *testing.T) {
			assert.True(t, isIPv4(item), "incorrect fail")
		})
	}
	for _, item := range fail {
		t.Run(item.String(), func(t *testing.T) {
			assert.False(t, isIPv4(item), "incorrect pass")
		})
	}
}

func TestIsIPv6(t *testing.T) {
	t.Parallel()
	pass := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::88"),
	}
	fail := []net.IP{
		net.ParseIP("192.0.2.1"),
		net.ParseIP("198.51.100.1"),
		net.ParseIP("203.0.113.1"),
		net.ParseIP("::ffff:192.0.2.2"),
		[]byte{1, 2, 3, 4, 5},
	}
	for _, item := range pass {
		t.Run(item.String(), func(t *testing.T) {
			assert.True(t, isIPv6(item), "incorrect fail")
		})
	}
	for _, item := range fail {
		t.Run(item.String(), func(t *testing.T) {
			assert.False(t, isIPv6(item), "incorrect pass")
		})
	}
}
