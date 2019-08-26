package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 参考情報 https://github.com/gregjones/httpcache

const (
	self   = ":8080"
	target = "https://api.spotify.com"
)

var (
	cache = make(map[string][]byte) // Cache Key:URL Value:response.Body
)

func main() {
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &myTransport{}
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(self, nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}

type myTransport struct {
}

func (t *myTransport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := http.DefaultTransport.(canceler); ok {
		cr.CancelRequest(req)
	}
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	//log.Print("===================REQUEST==========================")
	//log.Println(request.URL)
	//log.Printf("%#v\n", request)
	//log.Print("===================REQUEST==========================")

	cacheKey := request.URL.String()
	var canCache bool = false
	if request.Method == "GET" {
		canCache = true
	}
	if canCache {
		// get cache
		if cachedBody, ok := cache[cacheKey]; ok {
			log.Println("cache hit [" + request.Method + "] " + cacheKey)

			b := bytes.NewBuffer(cachedBody)
			response, err := http.ReadResponse(bufio.NewReader(b), request)
			if err != nil {
				return nil, err
			}
			return response, nil
		}
	}

	// NO CACHE
	log.Println("cache no  [" + request.Method + "] " + cacheKey)
	response, err := http.DefaultTransport.RoundTrip(request)
	//log.Print("===================RESPONSE==========================")
	//log.Printf("%s\n", string(body))
	//log.Print(string(body))
	//log.Print("===================RESPONSE==========================")

	// copying the response body
	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, err
	}

	// set cache
	if canCache {
		cache[cacheKey] = body
	}

	return response, err
}
