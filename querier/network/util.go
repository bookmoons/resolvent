package network

import (
	"fmt"
	"net"

	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func constructHostport(
	address net.IP,
	port uint16,
) (hostport string, err error) {
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

func constructClients() (
	clients map[string]map[resolvent.Protocol]*dns.Client,
	err error,
) {
	clients = make(map[string]map[resolvent.Protocol]*dns.Client)
	clients[net.IPv4zero.String()] = constructDefaultAddressClients()
	clients[net.IPv6zero.String()] = clients[net.IPv4zero.String()]
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	var ip net.IP
	for _, address := range addresses {
		ip, _, err = net.ParseCIDR(address.String())
		if err != nil {
			return
		}
		clients[ip.String()] = constructAddressClients(ip)
	}
	return
}

func constructDefaultAddressClients() map[resolvent.Protocol]*dns.Client {
	clients := make(map[resolvent.Protocol]*dns.Client)
	clients[resolvent.UDP] = &dns.Client{
		Net: "udp",
	}
	clients[resolvent.TCP] = &dns.Client{
		Net: "tcp",
	}
	return clients
}

func constructAddressClients(
	address net.IP,
) (clients map[resolvent.Protocol]*dns.Client) {
	clients = make(map[resolvent.Protocol]*dns.Client)
	clients[resolvent.UDP] = &dns.Client{
		Net: "udp",
		Dialer: &net.Dialer{
			LocalAddr: &net.UDPAddr{
				IP: address,
			},
		},
	}
	clients[resolvent.TCP] = &dns.Client{
		Net: "tcp",
		Dialer: &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP: address,
			},
		},
	}
	return
}
