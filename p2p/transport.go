package p2p

// Peer is an interface that represents the remot node
type Peer interface {
}

// Transport is anything that handles the
// between the nodes and the network.
// This can be {TCP,UDP,websockets, ...}
type Transport interface {
	listenAndAccept() error
	//This will listen for anything and everything!!!!
}
