package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

func main() {
	// Read filename from args
	fn := os.Args[1]
	log.Printf("Torrenting %v\n", fn)

	// Parse Torrent file
	r, err := os.Open(fn)
	if err != nil {
		log.Fatal("Coudn't open file ", err)
	}
	defer r.Close()

	var torrent Torrent

	err = bencode.Unmarshal(r, &torrent)
	if err != nil {
		log.Fatal("Unable to decode .torrent")
	}

	// Do torrent protocol
	fmt.Println(torrent.Pretty())
}

func (t *Torrent) Pretty() string {
	bs, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Println("Error marshalling", err)
		return ""
	}

	return string(bs)
}

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
