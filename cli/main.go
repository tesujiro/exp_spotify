package main

import (
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Print(`	cli search album|artist|playlist|track [keyword]
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
		"search": "https://api.spotify.com/v1/search",
	}
	switch cmd {
	case "search":
		search(token, endpoint[cmd], args)
	default:
		usage()
		os.Exit(1)
	}
}
