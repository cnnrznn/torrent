package client

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/cnnrznn/torrent/file"
	"github.com/google/uuid"
)

type Client struct {
	torrent file.Torrent
	peerID  string
	port    int
	peers   map[string]Peer

	sync.Mutex
}

func New(t file.Torrent) *Client {
	return &Client{
		torrent: t,
		peerID:  uuid.New().String()[:20],
		port:    6883,
		peers:   map[string]Peer{},
	}
}

func (c *Client) Run() {
	tChan := make(chan TrackerResponse)
	go c.PingTracker(tChan)

	for {
		select {
		case tr := <-tChan:
			fmt.Printf("%+v\n", tr)
			c.updatePeers(tr)
		}
	}
}

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

	/*am_choking := 1
	am_interested := 0
	peer_choking := 1
	peer_interested := 0*/

	// Connect to peer
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", peer.IP, peer.Port))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	c.doHandshake(peer, conn)
}

func (c *Client) doHandshake(peer Peer, conn net.Conn) {
	bs := c.buildHandshake(peer)
	log.Println("Handshake: ", string(bs))

	size := make([]byte, 1)
	n, err := conn.Read(size)
	if err != nil || n != 1 {
		log.Println(err)
		return
	}

	total_size := int(size[0]) + 48
	bs = make([]byte, total_size)

	n, err = conn.Read(bs)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(bs))
}

func (c *Client) buildHandshake(peer Peer) []byte {
	bs := make([]byte, 68)

	bs[0] = 19
	copy(bs[1:], "BitTorrent protocol")
	copy(bs[28:], c.torrent.InfoHash)
	copy(bs[48:], []byte(c.peerID))

	return bs
}
