package network

import (
	"context"
	"net"

	"github.com/glestaris/ice-agent/ice"
	"github.com/vishvananda/netlink"
)

// Networks composes and returns a list of iCE instance network entries that
// describe the available network interface in the running machine.
func Networks(ctx context.Context) ([]ice.InstanceNetwork, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	res := []ice.InstanceNetwork{}
	for _, link := range links {
		addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
		if err != nil {
			return nil, err
		}
		if len(addrs) == 0 {
			continue
		}

		for _, addr := range addrs {
			res = append(res, ice.InstanceNetwork{
				Iface:           addr.Label,
				IPAddr:          addr.IP,
				BroadcastIPAddr: net.ParseIP("0.0.0.0"),
			})
		}
	}

	return res, nil
}
