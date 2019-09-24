package route

import (
	"fmt"
	"net"

	"github.com/loadimpact/resolvent/internal"
)

func constructKey(local net.IP, address net.IP) (key string, err error) {
	localString, err := internal.ConstructAddress(local)
	if err != nil {
		return
	}
	addressString, err := internal.ConstructAddress(address)
	if err != nil {
		return
	}
	return fmt.Sprintf("%s-%s", localString, addressString), nil
}
