package network

import "testing"

func TestNetworks(t *testing.T) {
	nets, err := Networks(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(nets) != 2 {
		t.Fatalf("Expected 2 networks, got %d", len(nets))
	}
	for i, net := range nets {
		if net.Iface == "" {
			t.Fatalf("Expected non-empty intrface name for network %d", i)
		}
		if net.IPAddr == nil {
			t.Fatalf("Expected non-nil ip address for network %d", i)
		}
	}
}
