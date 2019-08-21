package main

type PagingBase struct {
	Href     string
	Limit    int
	Offset   int
	Total    int
	Next     string
	Previous string
}

type Albums struct {
	Albums PagingAlbums
}

type PagingAlbums struct {
	PagingBase
	Items []Album
}

type Artists struct {
	Artists PagingArtists
}

type PagingArtists struct {
	PagingBase
	Items []Artist
}

type ExternalURLs map[string]string

type Followers struct {
	Href  string
	Total int
}

type Image struct {
	Height int
	URL    string
	Width  int
}

type Restrictions map[string]string

type Album struct {
	AlbumGroup           string
	AlbumType            string
	Artists              []Artist
	AvailableMarkets     []string
	ExternalURLs         ExternalURLs
	Href                 string
	Id                   string
	Images               []Image
	Name                 string
	ReleaseDate          string
	ReleaseDatePrecision string
	Restrictions         Restrictions
	Type                 string
	URI                  string
}

type Artist struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	Followers    Followers
	Genres       []string
	Href         string
	Id           string
	Images       []Image
	Name         string
	Popularity   int
	Type         string
	URI          string
}
