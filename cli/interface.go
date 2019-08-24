package main

import "time"

type Timestamp time.Time

type PagingBase struct {
	Href     string
	Limit    int
	Offset   int
	Total    int
	Next     string
	Previous string
}

type Albums struct { // for Search
	Albums PagingAlbums
}

type PagingAlbums struct {
	PagingBase
	Items []Album
}

//type Artists struct { // for Search
//Artists PagingArtists
//}

//type PagingArtists struct {
//PagingBase
//Items []Artist
//}

//type Playlists struct { // for Search
//Playlists PagingPlaylists
//}

//type PagingPlaylists struct {
//PagingBase
//Items []Playlist
//}

type PagingPlaylistTracks struct {
	PagingBase
	Items []PlaylistTrack
}

type ExternalIDs map[string]string
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

type User struct {
	DisplayName  string
	ExternalURLs ExternalURLs
	Followers    Followers
	Href         string
	Id           string
	Images       []Image
	Type         string
	URI          string
}

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

type Playlist struct {
	Collaborative bool
	Description   string
	ExternalURLs  ExternalURLs `json:"external_urls"`
	Followers     Followers
	Href          string
	Id            string
	Images        []Image
	Name          string
	Owner         User
	Public        bool
	SnapshotId    string
	Tracks        PagingPlaylistTracks
	Type          string
	URI           string
}

type PlaylistTrack struct {
	AddedAt Timestamp
	AddedBy User
	IsLocal bool
	Track   Track
}

type Track struct {
	Album            Album
	Artists          []Artist
	AvailableMarkets []string
	DiscNumber       int
	DurationMs       int
	Explicit         bool
	ExternalIDs      ExternalIDs
	ExternalURLs     ExternalURLs
	Href             string
	Id               string
	IsPlayable       bool
	LinkedFrom       TrackLink
	Name             string
	Popularity       int
	PreviewURL       string
	TrackNumber      int
	Type             string
	URI              string
}

type TrackLink struct {
	ExternalURLs ExternalURLs
	Href         string
	Id           string
	Type         string
	URI          string
}
