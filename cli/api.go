package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

/*
type SpotifyAPI struct {
	cmd      string
	target   string
	usage    string
	desc     string
	endpoint string
}
*/

func get(token, endpoint string, params url.Values) ([]byte, error) {
	return call(token, "GET", endpoint, params, nil)
}

func put(token string, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	return call(token, "PUT", endpoint, params, body)
}

func call(token, method, endpoint string, params url.Values, body io.Reader) ([]byte, error) {

	//fmt.Println("endpoint:", endpoint)
	baseUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	baseUrl.RawQuery = params.Encode() // Escape Query Parameters

	req, err := http.NewRequest(method, baseUrl.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var client *http.Client
	if os.Getenv("ProxyServer") == "" {
		client = http.DefaultClient
	} else {
		proxyURL, err := url.Parse(os.Getenv("ProxyServer"))
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("proxyURL=%v\n", proxyURL)
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyURL(proxyURL),
		}
		client = &http.Client{
			Transport: transport,
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("resp=%v\n", resp)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("response.Body=%#v\n", resp.Body)
		return nil, fmt.Errorf("bad response status code %d", resp.StatusCode)
	}
	//fmt.Println("response status code ", resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
