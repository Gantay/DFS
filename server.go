package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"log"
	"sync"

	"github.com/Gantay/DFS/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootStrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) broadCast(p *Payload) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

type Payload struct {
	Key  string
	Data []byte
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1.Store this file to disk.
	// 2.Broadcast this file to all known peers in the network.

	if err := s.store.Write(key, r); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil
	}

	p := &Payload{
		Key:  key,
		Data: buf.Bytes(),
	}

	fmt.Println(buf.Bytes())

	return s.broadCast(p)
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("Connected with remote %s", p)

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action.")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():

			var p Payload
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", p)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) bootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}

		fmt.Println("attempting to connect with remote:", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {

				log.Println("dial error:", err)

			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootStrapNodes) != 0 {

		s.bootStrapNetwork()
	}

	s.loop()

	return nil
}
