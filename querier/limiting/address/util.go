package address

import (
	"net"

	"github.com/loadimpact/resolvent/internal"
)

func constructKey(address net.IP) (key string, err error) {
	return internal.ConstructAddress(address)
}
