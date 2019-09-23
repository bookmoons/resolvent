package internal

import (
	"net"
)

// IsIPv4 checks whether an IP address is IPv4.
func IsIPv4(address net.IP) bool {
	return address.To4() != nil
}

// IsIPv6 checks whether an IP address is IPv6.
func IsIPv6(address net.IP) bool {
	return address.To16() != nil && address.To4() == nil
}
