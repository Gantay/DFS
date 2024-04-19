package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	//conn is the underlying connection of the peer.
	conn net.Conn

	//if we dial and retrieve a conn => bool == true
	//if we  accept and retrieve a conn => bool == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPtransport struct {
	listenAddress string
	listner       net.Listener
	shakeHands    HandshakeFunc
	decoder       Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// in handshake.go func NOPHandShakeFunc(any) error { return nil }

func NewTCPTransport(listenAddr string) *TCPtransport {
	return &TCPtransport{
		shakeHands:    NOPHandShakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPtransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil

}

func (t *TCPtransport) startAcceptLoop() {
	for {
		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}
		fmt.Printf("New incomming connection %+v\n", conn)
		go t.handleconn(conn)
	}
}

type Temp struct{}

func (t *TCPtransport) handleconn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {

	}

	// Read loop
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("tcp error %s\n", err)
			continue
		}
	}

}
