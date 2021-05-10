package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func (c *Client) updatePeers(res TrackerResponse) {
	c.Lock()
	defer c.Unlock()

	for _, peer := range res.Peers {
		if _, ok := c.peers[peer.ID]; !ok {
			c.peers[peer.ID] = peer
			go c.handlePeer(peer)
		}
	}
}

func (c *Client) handlePeer(peer Peer) {
	defer func() {
		c.Lock()
		defer c.Unlock()

		delete(c.peers, peer.ID)
	}()

	/*am_interested := 0
	peer_choking := 1*/

	// Connect to peer
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", peer.IP, peer.Port))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	err = c.doHandshake(peer, conn)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Successful handshake with %+v\n", peer)

	// Tell peer I'm interested
	// Pull pieces from piece channel, try to download piece
}

func (c *Client) doHandshake(peer Peer, conn net.Conn) error {
	bs := c.buildHandshake(peer)

	n, err := conn.Write(bs)
	if err != nil || n != len(bs) {
		return err
	}

	size := make([]byte, 1)
	n, err = conn.Read(size)
	if err != nil || n != 1 {
		return err
	}

	total_size := int(size[0]) + 48
	bs = make([]byte, total_size)

	n, err = conn.Read(bs)
	if err != nil || n != len(bs) {
		return err
	}

	if bytes.Compare(c.torrent.InfoHash, bs[len(bs)-40:len(bs)-20]) != 0 {
		return fmt.Errorf("Info hash doesn't match during handshake")
	}

	//if bytes.Compare([]byte(peer.ID), bs[len(bs)-20:]) != 0 {
	//	return fmt.Errorf("Peer ID does not match the one provided by the tracker")
	//}

	return nil
}

func (c *Client) buildHandshake(peer Peer) []byte {
	bs := make([]byte, 68)

	bs[0] = 19
	copy(bs[1:], "BitTorrent protocol")
	copy(bs[28:], c.torrent.InfoHash)
	copy(bs[48:], []byte(c.peerID))

	return bs
}
