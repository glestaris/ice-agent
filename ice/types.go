package ice

import "net"

// InstanceNetwork is a part of the Instance type. It describes a network
// entry.
type InstanceNetwork struct {
	IPAddr          net.IP `json:"addr"`
	Iface           string `json:"iface"`
	BroadcastIPAddr net.IP `json:"bcast_addr"`
}

// Instance describes an iCE instance.
type Instance struct {
	ID                       string            `json:"id,omitempty"`
	SessionID                string            `json:"session_id"`
	Networks                 []InstanceNetwork `json:"networks"`
	PublicIPAddr             net.IP            `json:"public_ip_addr"`
	PublicReverseDNS         string            `json:"public_reverse_dns"`
	SSHUsername              string            `json:"ssh_username"`
	SSHAuthorizedFingerprint string            `json:"ssh_authorized_fingerprint"`
}
