package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Print(`	cli search album|artist|playlist|track [keyword]
	cli list devices
`)
}

func openbrowser(rawurl string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", rawurl).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawurl).Start()
	case "darwin":
		err = exec.Command("open", rawurl).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func getAccessToken() (string, error) {
	l, err := net.Listen("tcp", "localhost:8989")
	if err != nil {
		return "", err
	}
	defer l.Close()

	clientID := os.Getenv("ClientID")
	clientSecret := os.Getenv("ClientSecret")
	if clientID == "" || clientSecret == "" {
		err := fmt.Errorf("Env \"ClientID\", \"ClientSecret\" is not set")
		log.Fatal(err)
		return "", err
	}
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			// CAUTION: set scopes for APIs
			"user-read-playback-state",
			"playlist-read-private",
			"user-modify-playback-state",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: "http://localhost:8989", // CAUTION: this URL must be set on the Spotify dashboard
	}

	stateBytes := make([]byte, 16)
	_, err = rand.Read(stateBytes)
	if err != nil {
	}

	state := fmt.Sprintf("%x", stateBytes)
	//rawurl := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	rawurl := conf.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "token"))
	//fmt.Println("URL:", rawurl)

	// open in browser
	err = openbrowser(rawurl)
	if err != nil {
		return "", err
	}

	// Get Access token
	// see https://mattn.kaoriya.net/software/lang/go/20161231001721.htm
	// see https://qiita.com/TakahikoKawasaki/items/8567c80528da43c7e844#%E3%83%95%E3%83%A9%E3%82%B0%E3%83%A1%E3%83%B3%E3%83%88%E9%83%A8%E3%81%AF-http-%E3%83%AA%E3%82%AF%E3%82%A8%E3%82%B9%E3%83%88%E3%81%AB%E5%90%AB%E3%81%BE%E3%82%8C%E3%81%AA%E3%81%84
	quit := make(chan string)
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			//fmt.Println("HandlerFunc1!!")
			w.Write([]byte(`<script>location.href = "/close?" + location.hash.substring(1);</script>`))
		} else {
			//fmt.Println("HandlerFunc2!!")
			w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
			w.(http.Flusher).Flush()
			quit <- req.URL.Query().Get("access_token")
		}
	}))

	return <-quit, nil
}

func get(token string, endpoint string, query map[string]string) ([]byte, error) {
	baseUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	// set query parameters
	params := url.Values{}
	for k, v := range query {
		params.Add(k, v)
	}

	baseUrl.RawQuery = params.Encode() // Escape Query Parameters

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad response status code %d", resp.StatusCode)
	}
	//fmt.Println("response status code ", resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

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
