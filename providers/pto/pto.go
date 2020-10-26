package pto

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"crypto/tls"
	"time"

	out "github.com/silentmurdock/wrserver/providers/output"
)

type apiMovieResponse struct {
	Items []struct {
		Id string `json:"id"`
		Quality string `json:"quality"`
		SizeBytes int64 `json:"size_bytes"`
		Lang string `json:"language"`
		File string `json:"file"`
		TorrentSeeds int64 `json:"torrent_seeds"`
		TorrentPeers int64 `json:"torrent_peers"`
	} `json:"items"`
	ItemsLang []struct {
		Id string `json:"id"`
		Quality string `json:"quality"`
		SizeBytes int64 `json:"size_bytes"`
		Lang string `json:"language"`
		File string `json:"file"`
		TorrentSeeds int64 `json:"torrent_seeds"`
		TorrentPeers int64 `json:"torrent_peers"`
	} `json:"items_lang"`
}

func GetMovieMagnetByImdb(imdb string, ch chan<-[]out.OutputMovieStruct) {
	req, err := http.NewRequest("GET", "https://api.apiumadomain.com/movie?cb=&quality=720p,1080p&page=1&imdb=" + imdb, nil)
	if err != nil {
		ch <- []out.OutputMovieStruct{}
		return
	}

	//req.Header.Set("User-Agent", UserAgent)	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ch <- []out.OutputMovieStruct{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- []out.OutputMovieStruct{}
		return
	}

	response := apiMovieResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		ch <- []out.OutputMovieStruct{}
		return
	}

	outputMovieData := []out.OutputMovieStruct{}

	for _, thistorrent := range response.Items {
		if thistorrent.Quality != "3D" && thistorrent.SizeBytes <= 5368709120 {
			temp := out.OutputMovieStruct {
			    Hash: thistorrent.Id,
			    Quality: thistorrent.Quality,
			    Size: strconv.FormatInt(thistorrent.SizeBytes, 10),
			    Provider: "PTO",
			    Lang: thistorrent.Lang,
			    Title: out.RemoveFileExtension(thistorrent.File),
			    Seeds: strconv.FormatInt(thistorrent.TorrentSeeds, 10),
			    Peers: strconv.FormatInt(thistorrent.TorrentPeers, 10),
			}
			outputMovieData = append(outputMovieData, temp)
		}
	}

	for _, thistorrent := range response.ItemsLang {
		if thistorrent.Quality != "3D" && thistorrent.SizeBytes <= 5368709120 {
			temp := out.OutputMovieStruct {
			    Hash: thistorrent.Id,
			    Quality: thistorrent.Quality,
			    Size: strconv.FormatInt(thistorrent.SizeBytes, 10),
			    Provider: "PTO",
			    Lang: thistorrent.Lang,
			    Title: out.RemoveFileExtension(thistorrent.File),
			    Seeds: strconv.FormatInt(thistorrent.TorrentSeeds, 10),
			    Peers: strconv.FormatInt(thistorrent.TorrentPeers, 10),
			}
			outputMovieData = append(outputMovieData, temp)
		}
	}

	ch <- outputMovieData
}