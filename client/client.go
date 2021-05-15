package client

import (
	"log"
	"sync"

	"github.com/cnnrznn/torrent/file"
)

type Client struct {
	torrent file.Torrent
	peerID  string
	port    int
	peers   map[string]Peer

	sync.Mutex
}

func New(t file.Torrent) *Client {
	rs := make([]rune, 20)
	copy(rs, []rune("arandompeerid"))
	return &Client{
		torrent: t,
		//peerID:  uuid.New().String()[:20],
		peerID: string(rs),
		port:   6883,
		peers:  map[string]Peer{},
	}
}

func (c *Client) Run() {
	tChan := make(chan TrackerResponse, len(c.torrent.Announcers))
	for _, tracker := range c.torrent.Announcers {
		go c.PingTracker(tracker, tChan)
	}

	for {
		select {
		case tr := <-tChan:
			log.Printf("%+v\n", tr)
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
			go peer.Run()
		}
	}
}
