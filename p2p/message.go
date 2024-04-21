package p2p

import "net"

// Message holds any arbitrary data that is sent over the
// each transport between two node ion the network.
type Message struct {
	From    net.Addr
	Payload []byte
}
