package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Print(`	cli search album(s)|artist(s)|playlist(s)|track(s) [keyword]
	cli get profile [user_id]
	cli get playlist [playlist_id]
	cli get playlists [user_id]
	cli list device(s)
	cli list playlist(s)
	cli list profile
	cli play [device_id]
`)
}

func main() {
	var err error
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	token, err := getAccessToken() // CAUTION: The access tokens expire after 1 hour.
	if err != nil {
		log.Fatal("faild to get access token:", err)
	}
	//fmt.Println("token:", token)

	cmd := os.Args[1]
	args := os.Args[2:]

	base_url := "https://api.spotify.com"
	if os.Getenv("ReverseProxy") != "" {
		base_url = "http://localhost:8080"
	}
	endpoint := map[string]string{
		"devices/me": base_url + "/v1/me/player/devices",
		"search":     base_url + "/v1/search",
		//"player/me":    base_url+"/v1/me/player/",
		"play/me":      base_url + "/v1/me/player/play",
		"playlist":     base_url + "/v1/playlists/{playlist_id}",
		"playlists/me": base_url + "/v1/me/playlists",
		"playlists":    base_url + "/v1/users/{user_id}/playlists",
		//"profile/me":   base_url+"/v1/me",
		"profile/me": base_url + "/v1/me",
		"profile":    base_url + "/v1/users/{user_id}",
	}
	switch cmd {
	case "search":
		search(token, endpoint["search"], args)
	case "get":
		obj := args[0]
		args = args[1:]
		switch obj {
		case "profile":
			switch len(args) {
			case 0:
				profile(token, endpoint["profile/me"])
			case 1:
				ep := strings.ReplaceAll(endpoint["profile"], "{user_id}", args[0])
				profile(token, ep)
			default:
				usage()
				os.Exit(1)
			}
		case "playlist":
			switch len(args) {
			case 1:
				ep := strings.ReplaceAll(endpoint["playlist"], "{playlist_id}", args[0])
				playlist(token, ep)
			default:
				usage()
				os.Exit(1)
			}
		case "playlists":
			switch len(args) {
			case 0:
				playlists(token, endpoint["playlists/me"])
			case 1:
				ep := strings.ReplaceAll(endpoint["playlists"], "{user_id}", args[0])
				playlists(token, ep)
			default:
				usage()
				os.Exit(1)
			}
		default:
			usage()
			os.Exit(1)
		}
	case "create":
		obj := args[0]
		args = args[1:]
		switch obj {
		case "playlist":
		}
	case "list":
		if len(args) > 1 {
			usage()
			os.Exit(1)
		}
		obj := args[0]
		switch obj {
		case "device", "devices":
			devices(token, endpoint["devices/me"])
		case "playlists", "playlist":
			playlists(token, endpoint["playlists/me"])
		case "profile":
			profile(token, endpoint["profile/me"])
		default:
			usage()
			os.Exit(1)
		}
	case "play":
		switch len(args) {
		case 0:
			play(token, endpoint["play/me"], "")
		case 1:
			play(token, endpoint["play/me"], args[0])
		default:
			usage()
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}
