package client

import (
	"fmt"

	"github.com/cnnrznn/torrent/file"
	"github.com/google/uuid"
)

type Client struct {
	torrent file.Torrent
	peerID  string
	port    int
}

func New(t file.Torrent) *Client {
	return &Client{
		torrent: t,
		peerID:  uuid.New().String()[:20],
		port:    6883,
	}
}

func (c *Client) Run() {
	tChan := make(chan TrackerResponse)
	c.PingTracker(tChan)

	for {
		select {
		case tu := <-tChan:
			fmt.Println(tu)
		}
	}
}
