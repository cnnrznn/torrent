package main

import (
	"encoding/json"
	"log"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type Torrent struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Info         Info       `bencode:"info"`
}

type Info struct {
	Files       []File `bencode:"files"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type File struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func (t *Torrent) Pretty() string {
	bs, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Println("Error marshalling", err)
		return ""
	}

	return string(bs)
}
func ReadTorrent(fn string) (*Torrent, error) {
	// Parse Torrent file
	r, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var torrent Torrent

	err = bencode.Unmarshal(r, &torrent)
	if err != nil {
		return nil, err
	}

	return &torrent, nil
}
