package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/jackpal/bencode-go"
)

type TrackerRequest struct {
	InfoHash   string `url:"info_hash"`
	PeerID     string `url:"peer_id"`
	Port       int    `url:"port"`
	Uploaded   int    `url:"uploaded"`
	Downloaded int    `url:"downloaded"`
	Left       int    `url:"left"`
	Compact    int    `url:"compact"`
	NoPeerID   int    `url:"no_peer_id"`
	Event      string `url:"event"`
}

type TrackerResponse struct {
	FailureReason  string `bencode:"failure reason"`
	WarningMessage string `bencode:"warning message"`
	Interval       int    `bencode:"interval"`
	MinInterval    int    `bencode:"min interval"`
	TrackerID      string `bencode:"tracker id"`
	Complete       int    `bencode:"complete"`
	Peers          []Peer `bencode:"peers"`
}

func (c *Client) PingTracker(url string, ch chan<- TrackerResponse) {
	downloaded := 0
	uploaded := 0
	left := c.torrent.Size
	event := "started"
	interval := 5

	for ; ; time.Sleep(time.Duration(interval) * time.Second) {
		req := TrackerRequest{
			InfoHash:   string(c.torrent.InfoHash),
			PeerID:     c.peerID,
			Port:       c.port,
			Downloaded: downloaded,
			Uploaded:   uploaded,
			Left:       left,
			Event:      event,
		}

		res, err := SendPing(url, req)
		if err != nil {
			log.Println("Stopping tracker: ", err)
			return
		}

		if res.Interval > 0 {
			interval = res.Interval
		}
		if res.MinInterval > 0 {
			interval = res.MinInterval
		}

		event = ""

		ch <- *res
	}
}

func SendPing(url string, tReq TrackerRequest) (*TrackerResponse, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	vals, err := query.Values(tReq)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = vals.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(bs)

	var tres TrackerResponse
	err = bencode.Unmarshal(buf, &tres)
	if err != nil {
		return nil, err
	}
	return &tres, nil
}
