package p2p

import "net"

// Message holds any arbitrary data that is sent over the
// each transport between two node ion the network.
type RPC struct {
	From    net.Addr
	Payload []byte
}
