package network

import (
	"net"

	"github.com/loadimpact/resolvent"
	"github.com/miekg/dns"
)

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
