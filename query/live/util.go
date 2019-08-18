package live

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

func constructServer(
	address net.IP,
	port uint16,
) (server string, err error) {
	if isIPv6(address) {
		return fmt.Sprintf("[%s]:%d", address.String(), port), nil
	}
	if isIPv4(address) {
		return fmt.Sprintf("%s:%d", address.String(), port), nil
	}
	return "", errors.New("invalid IP address")
}

func isIPv4(address net.IP) bool {
	return address.To4() != nil
}

func isIPv6(address net.IP) bool {
	return address.To16() != nil && address.To4() == nil
}
