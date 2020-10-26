package itorrent

import (
	"net/http"
	"strconv"
	"regexp"
	"crypto/tls"
    "strings"
    "time"

	out "github.com/silentmurdock/wrserver/providers/output"

    "github.com/PuerkitoBio/goquery"
)

func GetMovieMagnetByImdb(imdb string, ch chan<-[]out.OutputMovieStruct) {
	// Disable security
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 10 * time.Second

	doc, err := goquery.NewDocument("https://itorrent.ws/torrentek/category/3/title/" + imdb + "/view_mode/photos/")
    if err != nil {
        ch <- []out.OutputMovieStruct{}
        return
    }

	outputMovieData := []out.OutputMovieStruct{}

	innerCh := make(chan out.OutputMovieStruct)

	counter := 0
    doc.Find("#ajaxtable .text-container").Each(func(_ int, item *goquery.Selection) {
        linkTag := item.Find("a")
        link, _ := linkTag.Attr("href")
        go scrapeMovieData(link, innerCh)
        counter++
    })

    for counter > 0 {
        temp := <-innerCh
        if (temp.Hash != "") {
        	outputMovieData = append(outputMovieData, temp)
        }
        counter--
    }

	ch <- outputMovieData	
}

func scrapeMovieData(movieUrl string, innerCh chan<-out.OutputMovieStruct) {
    doc, err := goquery.NewDocument("https://itorrent.ws" + movieUrl)
    if err != nil {
        innerCh <- out.OutputMovieStruct{}
        return
    }

    // Find title for raw magnet selection
    title := doc.Find("h1#torrent_title").Text()
    title = strings.TrimSpace(title)

    // Try to decode quality information from movieUrl
	quality := out.GuessQualityFromString(movieUrl)

    // Find Magnet link and decode infohash
    infoHash := ""
    doc.Find(".btn.btn-success.seed-warning").Each(func(_ int, item *goquery.Selection) {
        link, _ := item.Attr("href")
        infoHash = out.GetInfoHash(link)
    })
  
    size := ""
    language := ""
    seeds := ""
    leech := ""
    seedInt := int64(0)
    doc.Find("#torrent_page .left1").Each(func(_ int, item *goquery.Selection) {

        dataType := item.Find(".type").Text()
        switch dataType {
        case "Méret":
            size = out.DecodeSize(item.Next().Text())
        case "Peer":
        	value := item.Next().Text()
        	re := regexp.MustCompile("[0-9]+")
        	stringsize := re.FindAllString(value, -1)
        	seedInt, _ = strconv.ParseInt(stringsize[0], 10, 64)
            seeds = stringsize[0]
            leech = stringsize[1]
        case "Nyelv":
            language = out.DecodeLanguage(item.Next().Text(), "hu")
        }
    })

    if (seedInt == 0) {
    	innerCh <- out.OutputMovieStruct{}
        return
    }

    /*intSize, _ := strconv.ParseInt(size, 10, 64)
    if intSize > (5 * 1024 * 1024 * 1024) {
        innerCh <- out.OutputMovieStruct{}
    }*/

    innerCh <- out.OutputMovieStruct {
		    Hash: infoHash,
		    Quality: quality,
		    Size: size,
		    Provider: "ITORRENT",
		    Lang: language,
            Title: title,
            Seeds: seeds,
            Peers: leech,
	}
}

func GetShowMagnetByImdb(imdb string, season string, episode string, ch chan<-[]out.OutputShowStruct) {
    // Disable security
    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 10 * time.Second

    doc, err := goquery.NewDocument("https://itorrent.ws/torrentek/category/4/title/" + imdb + "/series_season/" + season + "/view_mode/photos/")
    if err != nil {
        ch <- []out.OutputShowStruct{}
        return
    }

    outputMovieData := []out.OutputShowStruct{}

    innerCh := make(chan out.OutputShowStruct)

    counter := 0
    doc.Find("#ajaxtable .text-container").Each(func(_ int, item *goquery.Selection) {
        linkTag := item.Find("a")
        link, _ := linkTag.Attr("href")
        go scrapeShowData(link, season, episode, innerCh)
        counter++
    })

    for counter > 0 {
        temp := <-innerCh
        if (temp.Hash != "") {
            outputMovieData = append(outputMovieData, temp)
        }
        counter--
    }

    ch <- outputMovieData   
}

func scrapeShowData(movieUrl string, season string, episode string, innerCh chan<-out.OutputShowStruct) {
    doc, err := goquery.NewDocument("https://itorrent.ws" + movieUrl)
    if err != nil {
        innerCh <- out.OutputShowStruct{}
        return
    }

    // Find title for raw magnet selection
    title := doc.Find("h1#torrent_title").Text()
    title = strings.TrimSpace(title)

    // Try to decode quality information from movieUrl
    quality := out.GuessQualityFromString(movieUrl)

    // Find Magnet link and decode infohash
    infoHash := ""
    doc.Find(".btn.btn-success.seed-warning").Each(func(_ int, item *goquery.Selection) {
        link, _ := item.Attr("href")
        infoHash = out.GetInfoHash(link)
    })
  
    size := ""
    language := ""
    seeds := ""
    leech := ""
    seedInt := int64(0)
    doc.Find("#torrent_page .left1").Each(func(_ int, item *goquery.Selection) {

        dataType := item.Find(".type").Text()
        switch dataType {
        case "Méret":
            size = out.DecodeSize(item.Next().Text())
        case "Peer":
            value := item.Next().Text()
            re := regexp.MustCompile("[0-9]+")
            stringsize := re.FindAllString(value, -1)
            seedInt, _ = strconv.ParseInt(stringsize[0], 10, 64)
            seeds = stringsize[0]
            leech = stringsize[1]
        case "Nyelv":
            language = out.DecodeLanguage(item.Next().Text(), "hu")
        }
    })

    if (seedInt == 0) {
        innerCh <- out.OutputShowStruct{}
        return
    }

    innerCh <- out.OutputShowStruct {
            Hash: infoHash,
            Quality: quality,
            Size: size,
            Provider: "ITORRENT",
            Lang: language,
            Title: title,
            Seeds: seeds,
            Peers: leech,
            Season: season,
            Episode: "0",
    }
}