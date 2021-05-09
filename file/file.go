package file

import (
	"crypto/sha1"
	"encoding/json"
	"log"
	"os"
	"strings"

	bencode "github.com/jackpal/bencode-go"
)

type Torrent struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Info         Info       `bencode:"info"`
	InfoHash     []byte
	Announcers   []string
	Size         int
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

	// Compute InfoHash
	h := sha1.New()
	err = bencode.Marshal(h, torrent.Info)
	if err != nil {
		return nil, err
	}
	torrent.InfoHash = h.Sum(nil)

	// TODO does this logic belong in 'client'?
	for _, row := range torrent.AnnounceList {
		for _, announcer := range row {
			if strings.HasPrefix(announcer, "http") {
				torrent.Announcers = append(torrent.Announcers, announcer)
			}
		}
	}

	for _, file := range torrent.Info.Files {
		torrent.Size += file.Length
	}

	return &torrent, nil
}
