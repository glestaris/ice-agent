package network

import (
	"context"
	"net"
	"strings"
)

// ReverseDNS queries the DNS server to get a FQDN for the provided IP address.
func ReverseDNS(ctx context.Context, ipAddr string) (string, error) {
	names, err := net.LookupAddr(ipAddr)
	if err != nil {
		return "", err
	}

	hostname := strings.Trim(names[0], ".")
	return hostname, nil
}
