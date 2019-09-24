package internal

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// ConstructAddress constructs an IP address string.
func ConstructAddress(address net.IP) (result string, err error) {
	if IsIPv6(address) || IsIPv4(address) {
		return address.String(), nil
	}
	return "", errors.New("invalid IP address")
}

// ConstructHostport constructs a hostport string from address and port.
func ConstructHostport(
	address net.IP,
	port uint16,
) (hostport string, err error) {
	if IsIPv6(address) {
		return fmt.Sprintf("[%s]:%d", address.String(), port), nil
	}
	if IsIPv4(address) {
		return fmt.Sprintf("%s:%d", address.String(), port), nil
	}
	return "", errors.New("invalid IP address")
}
