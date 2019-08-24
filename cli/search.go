package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func _search(token, endpoint, target string, args []string) ([]byte, error) {
	query := make(map[string]string)
	query["type"] = target
	query["limit"] = "50"
	for _, arg := range args {
		query["q"] = arg
	}
	return get(token, endpoint, query)
}

func search(token string, endpoint string, args []string) {
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}
	switch args[0] {
	case "album", "albums":
		b, err := _search(token, endpoint, "album", args[1:])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		var albums struct {
			Albums PagingAlbums
		}
		err = json.Unmarshal(b, &albums)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Println("Total:", albums.Albums.Total)
		// get artists info concurrently for each albums
		artists := make([][]Artist, len(albums.Albums.Items))
		wg := sync.WaitGroup{}
		for i, album := range albums.Albums.Items {
			wg.Add(1)
			go func(src []Artist, dst *[]Artist) {
				defer wg.Done()
				for _, artist := range src {
					b, err := get(token, artist.Href, nil)
					if err != nil {
						log.Print(err)
					}
					var a Artist
					err = json.Unmarshal(b, &a)
					if err != nil {
						log.Print(err)
					}
					*dst = append(*dst, a)
				}
			}(album.Artists, &artists[i])
		}
		wg.Wait()
		// display album info
		for i, album := range albums.Albums.Items {
			fmt.Printf("Album[%v]:\t", i)
			fmt.Printf("%v\t", album.Id)
			fmt.Printf("name:%v\t", album.Name)
			//fmt.Printf("artists:%#v\t", album.Artists)
			fmt.Printf("artists:")
			for _, artist := range artists[i] {
				fmt.Printf(" %v", artist.Name)
			}
			fmt.Printf("\n")
		}
	case "artist", "artists":
		b, err := _search(token, endpoint, "artist", args[1:])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		var artists struct {
			Artists PagingArtists
		}
		err = json.Unmarshal(b, &artists)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		//spew.Dump(artists)
		fmt.Println("Total:", artists.Artists.Total)
		for i, artist := range artists.Artists.Items {
			fmt.Printf("Artists[%v]:\t", i)
			fmt.Printf("%v\t", artist.Id)
			fmt.Printf("name:%v\t", artist.Name)
			fmt.Printf("\n")
		}
	case "playlist", "playlists":
		b, err := _search(token, endpoint, "playlist", args[1:])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		var playlists struct {
			Playlists PagingPlaylists
		}
		err = json.Unmarshal(b, &playlists)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Println("Total:", playlists.Playlists.Total)
		for i, playlist := range playlists.Playlists.Items {
			fmt.Printf("Playlist[%v]:\t", i)
			fmt.Printf("%v\t", playlist.Id)
			fmt.Printf("tracks:%v\t", playlist.Tracks.Total)
			fmt.Printf("name:%v\t", playlist.Name)
			fmt.Printf("\n")
		}
	case "track", "tracks":
		b, err := _search(token, endpoint, "track", args[1:])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		//fmt.Println("b=", string(b))
		var tracks struct {
			Tracks struct {
				PagingBase
				Items []Track
			}
		}
		err = json.Unmarshal(b, &tracks)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Println("Total:", tracks.Tracks.Total)
		for i, track := range tracks.Tracks.Items {
			fmt.Printf("Track[%v]:\t", i)
			fmt.Printf("%v\t", track.Id)
			fmt.Printf("name:%v\t", track.Name)
			fmt.Printf("%v (", track.Name)
			sep := ""
			for _, a := range track.Artists {
				fmt.Printf("%v%v", sep, a.Name)
				sep = ", "
			}
			fmt.Printf(") album: \"%v\"", track.Album.Name)
			fmt.Printf("\n")
		}
	default:
		usage()
		os.Exit(1)
	}
}
