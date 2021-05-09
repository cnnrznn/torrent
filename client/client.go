package client

import (
	"fmt"

	"github.com/cnnrznn/torrent/file"
	"github.com/google/uuid"
)

type Client struct {
	torrent file.Torrent
	peerID  uuid.UUID
	port    int
}

func New(t file.Torrent) *Client {
	return &Client{
		torrent: t,
		peerID:  uuid.New(),
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
