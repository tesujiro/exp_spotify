package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Print(`	cli search album|artist|playlist|track [keyword]
	cli profile [user_id]
	cli list devices
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
	endpoint := map[string]string{
		"search":     "https://api.spotify.com/v1/search",
		"profile/me": "https://api.spotify.com/v1/me",
		"profile":    "https://api.spotify.com/v1/users/{user_id}",
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
		}
	default:
		usage()
		os.Exit(1)
	}
}
