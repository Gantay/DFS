package p2p

import (
	"net"
)

// Peer is an interface that represents the remot node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything that handles the
// between the nodes and the network.
// This can be {TCP,UDP,websockets, ...}
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	//This will listen for anything and everything!!!!
	Consume() <-chan RPC
	Close() error
}
