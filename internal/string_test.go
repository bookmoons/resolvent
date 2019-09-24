package internal

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstructAddress(t *testing.T) {
	t.Parallel()
	t.Run("invalid", func(t *testing.T) {
		var address net.IP = []byte{1, 2, 3, 4, 5}
		_, err := ConstructAddress(address)
		require.EqualError(t, err, "invalid IP address")
	})
	pass := []struct {
		address net.IP
		result  string
	}{
		{net.ParseIP("192.0.2.1"), "192.0.2.1"},
		{net.ParseIP("2001:db8::1"), "2001:db8::1"},
	}
	for _, item := range pass {
		t.Run(item.address.String(), func(t *testing.T) {
			result, err := ConstructAddress(item.address)
			require.NoError(t, err, "incorrect fail")
			require.Equal(t, result, item.result)
		})
	}
}

func TestConstructHostport(t *testing.T) {
	t.Parallel()
	t.Run("invalid", func(t *testing.T) {
		var address net.IP = []byte{1, 2, 3, 4, 5}
		_, err := ConstructHostport(address, 50000)
		require.EqualError(t, err, "invalid IP address")
	})
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
			hostport, err := ConstructHostport(item.address, item.port)
			require.NoError(t, err, "incorrect fail")
			require.Equal(t, hostport, item.result)
		})
	}
}
