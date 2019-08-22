package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func playlists(token string, endpoint string) {
	b, err := get(token, endpoint, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	var playlists PagingPlaylists
	err = json.Unmarshal(b, &playlists)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	fmt.Println("Total:", playlists.Total)
	for i, playlist := range playlists.Items {
		fmt.Printf("Playlist[%v]:\t", i)
		fmt.Printf("tracks:%v\t", playlist.Tracks.Total)
		fmt.Printf("name:%v\t", playlist.Name)
		fmt.Printf("desc:%v\t", playlist.Description)
		fmt.Printf("\n")
	}
}
