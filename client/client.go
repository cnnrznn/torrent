package client

import (
	"fmt"

	"github.com/cnnrznn/torrent/file"
)

type Client struct {
	torrent file.Torrent
}

func New(t file.Torrent) *Client {
	return &Client{
		torrent: t,
	}
}

func (c *Client) Run() {
	fmt.Println("Run!")
}
