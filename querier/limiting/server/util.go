package server

import (
	"net"

	"github.com/loadimpact/resolvent/internal"
)

func constructKey(
	address net.IP,
	port uint16,
) (key string, err error) {
	return internal.ConstructHostport(address, port)
}
