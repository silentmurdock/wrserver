package providers

import (
	"net/url"
	"strings"
	"sort"
	"strconv"

	"github.com/silentmurdock/wrserver/providers/yts"
	"github.com/silentmurdock/wrserver/providers/eztv"
	"github.com/silentmurdock/wrserver/providers/pto"
	"github.com/silentmurdock/wrserver/providers/tmdb"
	"github.com/silentmurdock/wrserver/providers/itorrent"
	"github.com/silentmurdock/wrserver/providers/rarbg"
	"github.com/silentmurdock/wrserver/providers/x1337x"
	"github.com/silentmurdock/wrserver/providers/tvmaze"
	out "github.com/silentmurdock/wrserver/providers/output"
)

func GetMovieMagnet(imdbid string, query string, sources []string) []out.OutputMovieStruct {
	outputMovieData := []out.OutputMovieStruct{}

	ch := make(chan []out.OutputMovieStruct)

	counter := 0
	if imdbid != "" {
		for _, source := range sources {
			switch strings.ToLower(source) {
		    case "yts":
		    	go yts.GetMovieMagnetByImdb(imdbid, ch)
		    	counter++
		    case "itorrent":
		        go itorrent.GetMovieMagnetByImdb(imdbid, ch)
		        counter++
		    case "pto":
		        go pto.GetMovieMagnetByImdb(imdbid, ch)
		        counter++
		    case "rarbg":
		        go rarbg.GetMovieMagnetByImdb(imdbid, ch)
		        counter++
		    }	    
		}
	}

	if query != "" {
		params, err := url.ParseQuery(query)
		if err == nil {
			for _, source := range sources {
				switch strings.ToLower(source) {
			    case "1337x":
			        go x1337x.GetMovieMagnetByQuery(params, ch)
			        counter++
			    }
			}
		}
	}

	for counter > 0 {
		temp := <-ch	    
	    if len(temp) > 0 {
	    	if len(outputMovieData) > 0 {
		    	for _, item := range temp {
		    		duplicate := false
		    		for i, output := range outputMovieData {
		    			if strings.ToLower(output.Hash) == strings.ToLower(item.Hash) {
		    				duplicate = true
		    				if outputMovieData[i].Size == "0" && item.Size != "0" {
		    					outputMovieData[i].Size = item.Size
		    					outputMovieData[i].Title = item.Title
		    				}
		    			}
		    		}

		    		if duplicate == false {
		    			outputMovieData = append(outputMovieData, item)
		    		}		    		
		    	}
		    } else {
		    	for _, item := range temp {
		    		outputMovieData = append(outputMovieData, item)
		    	}
		    }
	    }
	    counter--	    
	}

	// Sort by seeds in descending order
	sort.Slice(outputMovieData, func(i, j int) bool {
		si, _ := strconv.ParseInt(outputMovieData[i].Seeds, 10, 64)
		sj, _ := strconv.ParseInt(outputMovieData[j].Seeds, 10, 64)
		return si > sj
	})
	
	return outputMovieData
}

func GetShowMagnet(imdbid string, query string, season string, episode string, sources []string) []out.OutputShowStruct {
	outputShowData := []out.OutputShowStruct{}

	ch := make(chan []out.OutputShowStruct)

	counter := 0
	if imdbid != "" {
		for _, source := range sources {
			switch strings.ToLower(source) {
		    case "eztv":
		        go eztv.GetShowMagnetByImdb(imdbid, season, episode, ch)
		        counter++
		    case "itorrent":
		        go itorrent.GetShowMagnetByImdb(imdbid, season, episode, ch)
		        counter++
		    case "rarbg":
		    	go rarbg.GetShowMagnetByImdb(imdbid, season, episode, ch)
		    	counter++
		    }
		}
	}

	if query != "" {
		params, err := url.ParseQuery(query)
		if err == nil {
			for _, source := range sources {
				switch strings.ToLower(source) {
			    case "1337x":
			        go x1337x.GetShowMagnetByQuery(params, season, episode, ch)
			        counter++
			    }
			}
		}
	}

	for counter > 0 {
		temp := <-ch
	    if len(temp) > 0 {
	    	if len(outputShowData) > 0 {
		    	for _, item := range temp {
		    		duplicate := false
		    		for i, output := range outputShowData {
		    			if strings.ToLower(output.Hash) == strings.ToLower(item.Hash) {
		    				duplicate = true
		    				if outputShowData[i].Size == "0" && item.Size != "0" {
		    					outputShowData[i].Size = item.Size
		    					outputShowData[i].Title = item.Title
		    				}
		    			}
		    		}

		    		if duplicate == false {
		    			outputShowData = append(outputShowData, item)
		    		}		    		
		    	}
		    } else {
		    	for _, item := range temp {
		    		outputShowData = append(outputShowData, item)
		    	}
		    }
	    }
	    counter--
	}

	// Sort by seeds in descending order
	sort.Slice(outputShowData, func(i, j int) bool {
		si, _ := strconv.ParseInt(outputShowData[i].Seeds, 10, 64)
		sj, _ := strconv.ParseInt(outputShowData[j].Seeds, 10, 64)
		return si > sj
	})
	
	return outputShowData
}

func SetTMDBKey(tmdbKey string) {
	tmdb.SetTMDBKey(tmdbKey)
}

func MirrorTmdbDiscover(qtype string, genretype string, sort string, date string, lang string, cpage string) string {
	return tmdb.MirrorTmdbDiscover(qtype, genretype, sort, date, lang, cpage)
}

func MirrorTmdbSearch(qtype string, lang string, cpage string, typedtext string) string {
	return tmdb.MirrorTmdbSearch(qtype, lang, cpage, typedtext)
}

func MirrorTmdbInfo(qtype string, tmdbid string, lang string) string {
	return tmdb.MirrorTmdbInfo(qtype, tmdbid, lang)
}

func GetTvMazeEpisodes(tvdb string, imdb string) string {
	return tvmaze.GetTvMazeEpisodes(tvdb, imdb)
}