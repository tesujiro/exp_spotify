package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func search(token string, endpoint string, args []string) {
	query := make(map[string]string)
	typ := []string{"album", "artist", "playlist", "track"}
	rev_type := make(map[string]bool)
	for _, k := range typ {
		rev_type[k] = true
	}
	if len(args) < 2 || !rev_type[args[0]] {
		usage()
		os.Exit(1)
	}
	query["type"] = args[0]
	query["limit"] = "50"
	for _, arg := range args[1:] {
		query["q"] = arg
	}
	b, err := get(token, endpoint, query)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	switch query["type"] {
	case "album":
		var albums Albums
		err = json.Unmarshal(b, &albums)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Println("Total:", albums.Albums.Total)
		for i, album := range albums.Albums.Items {
			fmt.Printf("Album[%v]:\t", i)
			fmt.Printf("release:%v\t", album.ReleaseDatePrecision)
			fmt.Printf("name:%v\t", album.Name)
			//fmt.Printf("artists:%#v\t", album.Artists)
			fmt.Printf("artists:")
			for _, artist := range album.Artists {
				b, err := get(token, artist.Href, nil)
				if err != nil {
					log.Print(err)
				}
				var a Artist
				err = json.Unmarshal(b, &a)
				if err != nil {
					log.Print(err)
				}
				fmt.Printf(" %v", a.Name)
			}
			fmt.Printf("\n")
		}
	case "artist":
		var artists Artists
		err = json.Unmarshal(b, &artists)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		//spew.Dump(artists)
		fmt.Println("Total:", artists.Artists.Total)
		for i, artist := range artists.Artists.Items {
			fmt.Printf("Artists[%v]:\t", i)
			fmt.Printf("name:%v\t", artist.Name)
			fmt.Printf("\n")
		}
	case "playlist":
		var playlists Playlists
		err = json.Unmarshal(b, &playlists)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Println("Total:", playlists.Playlists.Total)
		for i, playlist := range playlists.Playlists.Items {
			fmt.Printf("Playlist[%v]:\t", i)
			fmt.Printf("tracks:%v\t", playlist.Tracks.Total)
			fmt.Printf("name:%v\t", playlist.Name)
			fmt.Printf("\n")
		}
	}
}
