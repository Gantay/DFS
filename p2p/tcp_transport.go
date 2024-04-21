package p2p

import (
	"fmt"
	"net"
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

// close implement the Peer interface.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPtransport struct {
	TCPTransportOps
	listner net.Listener
	rpcch   chan RPC
}

func NewTCPTransport(opts TCPTransportOps) *TCPtransport {
	return &TCPtransport{
		TCPTransportOps: opts,
		rpcch:           make(chan RPC),
	}
}

// consume implements the Transport interface, which will return read-only channel.
// for reading the incoming messages reseived from another peer in the network.
func (t *TCPtransport) Consume() <-chan RPC {
	return t.rpcch

}

func (t *TCPtransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddr)
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
		fmt.Printf("New incomming connection big-W: %+v\n", conn)
		go t.handleconn(conn)
	}
}

func (t *TCPtransport) handleconn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP hanshake error: %s\n", err)
		conn.Close()
		return

	}

	// Read loop
	// reading from conn
	rpc := RPC{}
	for {

		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("tcp error %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()

		t.rpcch <- rpc

	}
}
