package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

type Peer struct {
	ID   string `bencode:"peer id"`
	IP   string `bencode:"ip"`
	Port int    `bencode:"port"`

	am_chocking   bool
	am_interested bool
	choked        bool
	interested    bool

	conn net.Conn

	client *Client
}

func (c *Client) updatePeers(res TrackerResponse) {
	c.Lock()
	defer c.Unlock()

	for _, peer := range res.Peers {
		if _, ok := c.peers[peer.ID]; !ok {
			c.peers[peer.ID] = peer
			go peer.handle()
		}
	}
}

func (p *Peer) handle() {
	defer func() {
		p.client.Lock()
		defer p.client.Unlock()

		delete(p.client.peers, p.ID)
	}()

	p.am_chocking = true
	p.am_interested = true
	p.choked = true
	p.interested = false

	// Connect to peer
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", peer.IP, peer.Port))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	p.conn = conn

	err = p.doHandshake()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Successful handshake with %+v\n", peer)

	// Tell peer I'm interested
	c.sendInterested(conn)

	// Pull pieces from piece channel, try to download piece
	for {
		msg, err := c.recvMsg(conn)
		if err != nil {
			log.Println(err)
			break
		}
		c.handleMsg(*msg)
	}
}

func (p *Peer) doHandshake() error {
	bs := p.buildHandshake()

	n, err := p.conn.Write(bs)
	if err != nil || n != len(bs) {
		return err
	}

	size := make([]byte, 1)
	n, err = p.conn.Read(size)
	if err != nil || n != 1 {
		return err
	}

	total_size := int(size[0]) + 48
	bs = make([]byte, total_size)

	n, err = p.conn.Read(bs)
	if err != nil || n != len(bs) {
		return err
	}

	if bytes.Compare(p.client.torrent.InfoHash, bs[len(bs)-40:len(bs)-20]) != 0 {
		return fmt.Errorf("Info hash doesn't match during handshake")
	}

	//if bytes.Compare([]byte(peer.ID), bs[len(bs)-20:]) != 0 {
	//	return fmt.Errorf("Peer ID does not match the one provided by the tracker")
	//}

	return nil
}

func (p *Peer) buildHandshake() []byte {
	bs := make([]byte, 68)

	bs[0] = 19
	copy(bs[1:], "BitTorrent protocol")
	copy(bs[28:], p.client.torrent.InfoHash)
	copy(bs[48:], []byte(p.client.peerID))

	return bs
}
