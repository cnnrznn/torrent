package client

import (
	"fmt"
	"log"

	"github.com/cnnrznn/torrent/file"
)

type TrackerResponse struct {
	FailureReason  string `bencode:"failue reason"`
	WarningMessage string `bencode:"warning message"`
	Interval       int    `bencode:"interval"`
	MinInterval    int    `bencode:"min interval"`
	TrackerID      string `bencode:"tracker id"`
	Complete       int    `bencode:"complete"`
	Peers          []Peer `bencode:"peers"`
}

type Peer struct {
	ID   string `bencode:"peer id"`
	IP   string `bencode:"ip"`
	port int    `bencode:"port"`
}

type Client struct {
	torrent file.Torrent
}

func New(t file.Torrent) *Client {
	return &Client{
		torrent: t,
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

func (c *Client) PingTracker(ch chan<- TrackerResponse) {
	log.Println(c.torrent.Announcers)
}
