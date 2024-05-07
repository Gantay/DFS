package main

import "github.com/Gantay/DFS/p2p"

func main() {
	tcpTransportOpts := p2p.TCPTransportOps{}
	tcpTransport := p2p.NewTCPTransport()

	fileServerOpts := FileServerOpts{
		ListenAddr:        "3000",
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewFileServer()

}
