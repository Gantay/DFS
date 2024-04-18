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

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPtransport {
	return &TCPtransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPtransport) listenAndAccept() error {
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
		go t.handleconn(conn)
	}
}

func (t *TCPtransport) handleconn(conn net.Conn) {
	fmt.Printf("New incomming connection %+v\n", conn)
}
