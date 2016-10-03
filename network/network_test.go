package network

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/ice-stuff/ice-agent/ice"
	"github.com/vishvananda/netlink"
)

func TestNetworks(t *testing.T) {
	nets, err := Networks(nil)
	if err != nil {
		t.Fatal(err)
	}

	link, err := netlink.LinkByName("eth0")
	if err != nil {
		t.Fatal(err)
	}
	addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		t.Fatal(err)
	}

	expectedNets := []ice.InstanceNetwork{
		ice.InstanceNetwork{
			Iface:           "lo",
			IPAddr:          net.ParseIP("127.0.0.1"),
			BroadcastIPAddr: net.ParseIP("0.0.0.0"),
		},
		ice.InstanceNetwork{
			Iface:           "eth0",
			IPAddr:          addrs[0].IP,
			BroadcastIPAddr: net.ParseIP("0.0.0.0"),
		},
	}

	netsJSON, err := json.Marshal(nets)
	if err != nil {
		t.Fatal(err)
	}
	netsJSONStr := string(netsJSON)
	expectedNetsJSON, err := json.Marshal(expectedNets)
	if err != nil {
		t.Fatal(err)
	}
	expectedNetsJSONStr := string(expectedNetsJSON)
	if netsJSONStr != expectedNetsJSONStr {
		t.Fatalf(
			"Expected list of networks:\n%#v\ngot\n%#v",
			expectedNetsJSONStr, netsJSONStr,
		)
	}
}
