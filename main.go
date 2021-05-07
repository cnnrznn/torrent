package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Read filename from args
	fn := os.Args[1]
	log.Printf("Torrenting %v\n", fn)

	torrent, err := ReadTorrent(fn)
	if err != nil {
		log.Fatal("Probem reading torrent: ", err)
	}

	// Do torrent protocol
	fmt.Println(torrent.Pretty())
}
