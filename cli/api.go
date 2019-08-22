package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
