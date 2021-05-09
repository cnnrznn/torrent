package client

import (
	"log"
	"time"
)

type TrackerRequest struct {
	InfoHash   []byte `bencode:"info_hash"`
	PeerID     string `bencode:"peer_id"`
	Port       int    `bencode:"port"`
	Uploaded   int    `bencode:"uploaded"`
	Downloaded int    `bencode:"downloaded"`
	Left       int    `bencode:"left"`
	Compact    int    `bencode:"compact"`
	NoPeerID   int    `bencode:"no_peer_id"`
	Event      string `bencode:"event"`
}

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

func (c *Client) PingTracker(ch chan<- TrackerResponse) {
	downloaded := 0
	uploaded := 0
	left := 0
	event := "started"
	interval := 0

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		tr := TrackerRequest{
			InfoHash:   c.torrent.InfoHash,
			PeerID:     c.peerID.String(),
			Port:       c.port,
			Downloaded: downloaded,
			Uploaded:   uploaded,
			Left:       left,
			Event:      event,
		}
		log.Println(tr)
	}
}
