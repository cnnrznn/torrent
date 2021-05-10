package main

import (
	"log"
	"os"

	"github.com/cnnrznn/torrent/client"
	"github.com/cnnrznn/torrent/file"
)

func main() {
	// Read filename from args
	fn := os.Args[1]

	torrent, err := file.ReadTorrent(fn)
	if err != nil {
		log.Fatal("Probem reading torrent: ", err)
	}

	log.Printf("Torrenting %v\n", torrent.Info.Name)

	// Do torrent protocol
	client := client.New(*torrent)
	client.Run()
}
