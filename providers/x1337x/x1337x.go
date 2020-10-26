package x1337x

import (
	"net/http"
	"crypto/tls"
    "strings"
    "time"

	out "github.com/silentmurdock/wrserver/providers/output"

    "github.com/PuerkitoBio/goquery"
)

func GetMovieMagnetByQuery(params map[string][]string, ch chan<-[]out.OutputMovieStruct) {
	// Decode params data
    title := ""
    releaseYear := ""

    if params["title"] != nil {
        title = params["title"][0]
    }

    if params["releaseyear"] != nil {
        releaseYear = params["releaseyear"][0]
    }

    if title == "" || releaseYear == "" {
        ch <- []out.OutputMovieStruct{}
        return
    }

    // Disable security
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 10 * time.Second

	doc, err := goquery.NewDocument("https://www.1337x.to/category-search/" + title + " " +  releaseYear + "/Movies/1/")
    if err != nil {
        ch <- []out.OutputMovieStruct{}
        return
    }

	outputMovieData := []out.OutputMovieStruct{}

	innerCh := make(chan out.OutputMovieStruct)

	counter := 0
    doc.Find("tbody tr").Each(func(_ int, item *goquery.Selection) {
        seedsClass := item.Find("td.seeds")
        seeds := seedsClass.Text()

        if seeds != "0" {
            nameClass := item.Find("td.name")        
            linkTag := nameClass.Find("a").Next()
            link, _ := linkTag.Attr("href")

            go scrapeMovieData(link, innerCh)
            counter++
        }
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

func GetShowMagnetByQuery(params map[string][]string, season string, episode string, ch chan<-[]out.OutputShowStruct) {
    // Decode params data
    title := ""

    if params["title"] != nil {
        title = params["title"][0]
    }

    if title == "" || (season == "0" && episode == "0") {
        ch <- []out.OutputShowStruct{}
        return
    }

    query := strings.ReplaceAll(title, " ", "+") + "+"

    if len(season) == 1 {
        query = query + "s0" + season
    } else {
        query = query + "s" + season
    }

    if len(episode) == 1 {
        query = query + "e0" + episode
    } else {
        query = query + "e" + episode
    }

    // Disable security
    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 10 * time.Second

    doc, err := goquery.NewDocument("https://www.1337x.to/category-search/" + query + "/TV/1/")
    if err != nil {
        ch <- []out.OutputShowStruct{}
        return
    }

    outputShowData := []out.OutputShowStruct{}

    innerCh := make(chan out.OutputShowStruct)

    counter := 0
    doc.Find("tbody tr").Each(func(_ int, item *goquery.Selection) {
        seedsClass := item.Find("td.seeds")
        seeds := seedsClass.Text()

        if seeds != "0" {
            nameClass := item.Find("td.name")        
            linkTag := nameClass.Find("a").Next()
            link, _ := linkTag.Attr("href")

            go scrapeShowData(link, season, episode, innerCh)
            counter++
        }
    })

    for counter > 0 {
        temp := <-innerCh
        if (temp.Hash != "") {
            outputShowData = append(outputShowData, temp)
        }
        counter--
    }

    ch <- outputShowData   
}

func scrapeMovieData(movieUrl string, innerCh chan<-out.OutputMovieStruct) {
    doc, err := goquery.NewDocument("https://www.1337x.to" + movieUrl)
    if err != nil {
        innerCh <- out.OutputMovieStruct{}
    }

    // Find title for raw magnet selection
    title := doc.Find("title").Text()
    title = strings.TrimPrefix(title, "Download")
    title = strings.TrimSuffix(title, "Torrent | 1337x")
    title = out.CleanString(title)

    // Trimmed title
    //title := doc.Find(".box-info-heading h1").Text()
    //title = strings.TrimSpace(title)

    // Try to decode quality information from movieUrl
    quality := out.GuessQualityFromString(movieUrl)

    // Find Magnet link and decode infohash
    infoHash := ""
    doc.Find(".torrent-detail-page ul li a").Each(func(_ int, item *goquery.Selection) {
        if item.Text() == "Magnet Download" {
            link, _ := item.Attr("href")
            infoHash = out.GetInfoHash(link)
        }
    })
  
    size := ""
    language := ""
    seeders := ""
    leechers := ""
    doc.Find(".torrent-detail-page ul.list li").Each(func(_ int, item *goquery.Selection) {
        textNode := item.ChildrenFiltered("strong").Text()
        if textNode == "Total size" {
            size = out.DecodeSize(item.ChildrenFiltered("span").Text())
        } else if textNode == "Language" {
            language = out.DecodeLanguage(item.ChildrenFiltered("span").Text(), "en")
        } else if textNode == "Seeders" {
            seeders = item.ChildrenFiltered("span").Text()
        } else if textNode == "Leechers" {
            leechers = item.ChildrenFiltered("span").Text()
        }
    })

    innerCh <- out.OutputMovieStruct {
		    Hash: infoHash,
		    Quality: quality,
		    Size: size,
		    Provider: "1337X",
		    Lang: language,
            Title: title,
            Seeds: seeders,
            Peers: leechers,
	}
}

func scrapeShowData(movieUrl string, season string, episode string, innerCh chan<-out.OutputShowStruct) {
    doc, err := goquery.NewDocument("https://www.1337x.to" + movieUrl)
    if err != nil {
        innerCh <- out.OutputShowStruct{}
    }

    // Find title for raw magnet selection
    title := doc.Find("title").Text()
    title = strings.TrimPrefix(title, "Download")
    title = strings.TrimSuffix(title, "Torrent | 1337x")
    title = out.CleanString(title)

    // Trimmed title
    //title := doc.Find(".box-info-heading h1").Text()
    //title = strings.TrimSpace(title)

    // Try to decode quality information from movieUrl
    quality := out.GuessQualityFromString(movieUrl)

    // Find Magnet link and decode infohash
    infoHash := ""
    doc.Find(".torrent-detail-page ul li a").Each(func(_ int, item *goquery.Selection) {
        if item.Text() == "Magnet Download" {
            link, _ := item.Attr("href")
            infoHash = out.GetInfoHash(link)
        }
    })
  
    size := ""
    language := ""
    seeders := ""
    leechers := ""
    doc.Find(".torrent-detail-page ul.list li").Each(func(_ int, item *goquery.Selection) {
        textNode := item.ChildrenFiltered("strong").Text()
        if textNode == "Total size" {
            size = out.DecodeSize(item.ChildrenFiltered("span").Text())
        } else if textNode == "Language" {
            language = out.DecodeLanguage(item.ChildrenFiltered("span").Text(), "en")
        } else if textNode == "Seeders" {
            seeders = item.ChildrenFiltered("span").Text()
        } else if textNode == "Leechers" {
            leechers = item.ChildrenFiltered("span").Text()
        }
    })

    innerCh <- out.OutputShowStruct {
            Hash: infoHash,
            Quality: quality,
            Size: size,
            Provider: "1337X",
            Lang: language,
            Title: title,
            Seeds: seeders,
            Peers: leechers,
            Season: season, // Need to rewrite later to scrape from url
            Episode: episode,
    }
}