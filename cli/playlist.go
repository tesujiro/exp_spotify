package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func playlist(token string, endpoint string) {
	b, err := get(token, endpoint, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	//fmt.Printf("%+v\n", string(b))
	var playlist Playlist
	err = json.Unmarshal(b, &playlist)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	//fmt.Printf("Playlist: %+v\n", playlist)
	fmt.Printf("ID: %v\n", playlist.Id)
	fmt.Printf("Desc: %v\n", playlist.Description)
	fmt.Printf("Name: %v\n", playlist.Name)
	//fmt.Printf("Owner: %+v\n", playlist.Owner)
	fmt.Printf("Tracks: %v\n", playlist.Tracks.Total)
	for i, ptrack := range playlist.Tracks.Items {
		fmt.Printf("Track[%v]: %v (", i, ptrack.Track.Name)
		sep := ""
		for _, a := range ptrack.Track.Album.Artists {
			fmt.Printf("%v%v", sep, a.Name)
			sep = ", "
		}
		fmt.Printf(") Album: \"%v\"\n", ptrack.Track.Album.Name)
	}
}

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
		fmt.Printf("Id:%v\t", playlist.Id)
		fmt.Printf("tracks:%v\t", playlist.Tracks.Total)
		fmt.Printf("name:%v\t", playlist.Name)
		fmt.Printf("desc:%v\t", playlist.Description)
		fmt.Printf("\n")
	}
}
