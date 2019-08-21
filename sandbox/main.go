package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

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
	fmt.Println("URL:", rawurl)

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
			fmt.Println("HandlerFunc1!!")
			w.Write([]byte(`<script>location.href = "/close?" + location.hash.substring(1);</script>`))
		} else {
			fmt.Println("HandlerFunc2!!")
			w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
			w.(http.Flusher).Flush()
			quit <- req.URL.Query().Get("access_token")
		}
	}))

	return <-quit, nil
}

func searchArtist(token string, keyword string) error {

	endpoint := "https://api.spotify.com/v1/search"
	baseUrl, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	// set query parameters
	params := url.Values{}
	params.Add("q", keyword)
	params.Add("type", "artist")

	baseUrl.RawQuery = params.Encode() // Escape Query Parameters

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("bad response status code %d", resp.StatusCode)
	}
	fmt.Println("response status code ", resp.StatusCode)

	/*
		byteArray, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response body:", string(byteArray))
	*/

	var ret interface{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return err
	}
	fmt.Printf("Response body: %v\n", ret)

	return nil
}

func main() {
	var err error

	token, err := getAccessToken() // CAUTION: The access tokens expire after 1 hour.
	if err != nil {
		log.Fatal("faild to get access token:", err)
	}
	fmt.Println("token:", token)

	err = searchArtist(token, "吾妻")
	if err != nil {
		log.Fatalf("spotifyGet error: %v", err)
	}

}
