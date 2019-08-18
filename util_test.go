package resolvent

import (
	"net"
	"testing"
)

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
			if !isIPv4(item) {
				t.Errorf("incorrect fail")
			}
		})
	}
	for _, item := range fail {
		t.Run(item.String(), func(t *testing.T) {
			if isIPv4(item) {
				t.Errorf("incorrect pass")
			}
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
			if !isIPv6(item) {
				t.Errorf("incorrect fail")
			}
		})
	}
	for _, item := range fail {
		t.Run(item.String(), func(t *testing.T) {
			if isIPv6(item) {
				t.Errorf("incorrect pass")
			}
		})
	}
}
