package p2p

// HandshakeFunc ...?
type HandshakeFunc func(Peer) error

func NOPHandShakeFunc(Peer) error { return nil }
