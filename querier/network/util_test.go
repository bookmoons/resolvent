package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructHostport(t *testing.T) {
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
			hostport, err := constructHostport(item.address, item.port)
			assert.NoError(t, err, "incorrect fail")
			assert.Equal(t, hostport, item.result)
		})
	}
	t.Run("invalid", func(t *testing.T) {
		var address net.IP = []byte{1, 2, 3, 4, 5}
		_, err := constructHostport(address, 50000)
		assert.EqualError(t, err, "invalid IP address")
	})
}
