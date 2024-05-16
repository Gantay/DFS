package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	//The underlying connection of the peer, which in this case
	// is a TCP connection.
	net.Conn
	//if we dial and retrieve a conn => bool == true
	//if we  accept and retrieve a conn => bool == false
	outbound bool

	Wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err

}

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
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

// Close implements the Transport interface.
func (t *TCPtransport) Close() error {
	return t.listner.Close()
}

// Dial implements the Transport interface.
func (t *TCPtransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err

	}
	fmt.Println(conn)

	go t.handleconn(conn, true)

	return nil

}

func (t *TCPtransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port: %s", t.ListenAddr)

	return nil

}

func (t *TCPtransport) startAcceptLoop() {
	for {
		conn, err := t.listner.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}

		go t.handleconn(conn, false)
	}
}

func (t *TCPtransport) handleconn(conn net.Conn, outBound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outBound)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP hanshake error: %s\n", err)
		conn.Close()
		return

	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	// reading from conn
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr().String()
		peer.Wg.Add(1)
		fmt.Println("waiting till steam is done")
		t.rpcch <- rpc
		peer.Wg.Wait()
		fmt.Println("steam done continueing normal read loop")

	}
}
