package main

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/oz/osdb"

	"golang.org/x/text/encoding/charmap"

	"github.com/silentmurdock/wrserver/providers"
)

var (
	urlAPI = "/api/"

	// TMDB API key
	TMDBKey = "a4d9ad8d2d072c50dc998cc0d1a508fa"
	// OpenSubtitles user agent string
	OSUserAgent = "White Raven v0.3"
	
	upgrader = websocket.Upgrader{
	    ReadBufferSize:  1024,
	    WriteBufferSize: 1024,
	}
)

func setOSUserAgent(userAgent string) {
	OSUserAgent = userAgent
}

func setTMDBKey(tmdbKey string) {
	providers.SetTMDBKey(tmdbKey)
}

func fetchZip(zipurl string) (*zip.Reader, error) {
	req, err := http.NewRequest("GET", zipurl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", OSUserAgent)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(resp.Status)
		}
		return nil, errors.New(string(b))
	}

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(buf.Bytes())
	return zip.NewReader(b, int64(b.Len()))
}

func decodeData(encData []byte, enc string) string {
	dec := charmap.Windows1250.NewDecoder()
	switch enc {
	case "CP1251":
		dec = charmap.Windows1251.NewDecoder()
	case "CP1252":
		dec = charmap.Windows1252.NewDecoder()
	case "CP1253":
		dec = charmap.Windows1253.NewDecoder()
	case "CP1254":
		dec = charmap.Windows1254.NewDecoder()
	case "CP1255":
		dec = charmap.Windows1255.NewDecoder()
	case "CP1256":
		dec = charmap.Windows1256.NewDecoder()
	case "CP1257":
		dec = charmap.Windows1257.NewDecoder()
	case "CP1258":
		dec = charmap.Windows1258.NewDecoder()	
	}
    out, _ := dec.Bytes(encData)
    return string(out)
}

func handleAPI(cors bool) {
	routerAPI := mux.NewRouter()
	routerAPI.SkipClean(true)

	routerAPI.HandleFunc(urlAPI+"about", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, serverInfo())
	})

	routerAPI.HandleFunc(urlAPI+"stop", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, serverStop())
		if err == nil {
			go func() {
				time.Sleep(1 * time.Nanosecond)
				procQuit <- true
			}()
		} else {
			go func() {
				time.Sleep(1 * time.Nanosecond)
				procQuit <- true
			}()
		}
	})

	routerAPI.HandleFunc(urlAPI+"restart/downrate/{downrate}/uprate/{uprate}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		for _, thistorrent := range torrents {
			log.Println("Delete torrent:", thistorrent.torrent.InfoHash().String())
			stopAllFileDownload(thistorrent.torrent.Files())
			thistorrent.torrent.Drop()
			delete(torrents, thistorrent.torrent.InfoHash().String())		
		}
		
		_, err := io.WriteString(w, restartServer())
		if err == nil {
			go func() {
				time.Sleep(1 * time.Nanosecond)
				dr, _ := strconv.ParseInt(vars["downrate"], 10, 64)
				ur, _ := strconv.ParseInt(vars["uprate"], 10, 64)
				procRestart <- []int64 {dr, ur}
			}()
		} else {
			go func() {
				time.Sleep(1 * time.Nanosecond)
				dr, _ := strconv.ParseInt(vars["downrate"], 10, 64)
				ur, _ := strconv.ParseInt(vars["uprate"], 10, 64)
				procRestart <- []int64 {dr, ur}
			}()
		}
	})

	routerAPI.HandleFunc(urlAPI+"add/{hash}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		magnet := "magnet:?xt=urn:btih:" + vars["hash"]

		for tryCount := 0; tryCount < 4; tryCount++ {
			if tryCount > 0 {
				time.Sleep(10 * time.Second)
			}

			t := addMagnet(magnet)
			if t != nil {
				log.Println("Add torrent:", vars["hash"])
				io.WriteString(w, jsonFilesList(r.Host, t.Files()))
				return
			} else if len(torrents) == 0 {
				return
			}
		}

		if len(torrents) > 0 {
			io.WriteString(w, onlyOneTorrent())
		}
	})

	routerAPI.HandleFunc(urlAPI+"delete/{hash}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if t, ok := torrents[vars["hash"]]; ok {

			log.Println("Delete torrent:", vars["hash"])
			stopAllFileDownload(t.torrent.Files())
			t.torrent.Drop()
			delete(torrents, vars["hash"])

			io.WriteString(w, deleteTorrent())
		} else {
			log.Println("Already deleted torrent:", vars["hash"])
			return
		}
	})

	routerAPI.HandleFunc(urlAPI+"deleteall", func(w http.ResponseWriter, r *http.Request) {
		for _, thistorrent := range torrents {
			log.Println("Delete torrent:", thistorrent.torrent.InfoHash().String())
			stopAllFileDownload(thistorrent.torrent.Files())
			thistorrent.torrent.Drop()
			delete(torrents, thistorrent.torrent.InfoHash().String())		
		}
		io.WriteString(w, deleteAllTorrent())
	})

	routerAPI.HandleFunc(urlAPI+"get/{hash}/{base64path}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if d, err := base64.StdEncoding.DecodeString(vars["base64path"]); err == nil {
			if t, ok := torrents[vars["hash"]]; ok {

				idx := getFileByPath(string(d), t.torrent.Files())
				file := t.torrent.Files()[idx]

				path := file.DisplayPath()
				log.Println("Download torrent:", vars["hash"])

				incFileClients(path, t)

				runtime.GC()
				/*log.Println("Calculate Opensubtitles hash...")
				fileHash := calculateOpensubtitlesHash(file)
				log.Println("Opensubtitles hash calculated:", fileHash)*/

				serveTorrentFile(w, r, file)
				//stop downloading the file when no connections left
				if decFileClients(path, t) <= 0 {
					stopDownloadFile(file)					
				}
			} else {
				log.Println("Unknown torrent:", vars["hash"])
				return
			}
		} else {
			log.Println(err)
			return
		}
	})

	routerAPI.HandleFunc(urlAPI+"stats/{hash}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		
		if t, ok := torrents[vars["hash"]]; ok {					
			log.Println("Check torrent stats:", vars["hash"])
			io.WriteString(w, downloadStats(r.Host, t.torrent))
		} else {
			log.Println("Unknown torrent:", vars["hash"])
			return
		}
	})

	routerAPI.HandleFunc(urlAPI+"subtitlesbyimdb/{imdb}/lang/{lang}/season/{season}/episode/{episode}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// Create Opensubtitles client
		c, err := osdb.NewClient()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		c.UserAgent = OSUserAgent

		// Anonymous Login with UserAgent string will set c.Token when successful
		if err = c.LogIn("", "", ""); err != nil {
			http.NotFound(w, r)
			return
		}

		ids := make([]string, 1)
		ids[0] = strings.TrimPrefix(vars["imdb"], "tt")
		langs := strings.Split(vars["lang"], ",")
		
		// Fallback language always English
		fallbackLang := false;
		for _, l := range langs {
			if l == "eng" {
				fallbackLang = true
				break
			}
		}

		if fallbackLang == false {
			langs = append(langs, "eng")
		}

		log.Println("Search subtitle by imdbid...")
		
		season, err := strconv.ParseInt(vars["season"], 10, 64)
		episode, err :=	strconv.ParseInt(vars["episode"], 10, 64)

		params := []interface{}{}
		if season == 0 && episode == 0 {
			params = []interface{}{
				c.Token,
				[]struct {
					Imdb  string `xmlrpc:"imdbid"`
					Langs string `xmlrpc:"sublanguageid"`
				}{{
					ids[0],
					strings.Join(langs, ","),
				}},
			}
		} else {
			params = []interface{}{
				c.Token,
				[]struct {
					Imdb  string `xmlrpc:"imdbid"`
					Langs string `xmlrpc:"sublanguageid"`
					Season int64  `xmlrpc:"season"`
					Episode int64  `xmlrpc:"episode"`
				}{{
					ids[0],
					strings.Join(langs, ","),
					season,
					episode,
				}},
			}
		}

		res, err := c.SearchSubtitles(&params)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		log.Println("Subtitle found.")
		io.WriteString(w, subtitleFilesList(r.Host, res, langs[0]))
	})

	routerAPI.HandleFunc(urlAPI+"subtitlesbytext/{text}/lang/{lang}/season/{season}/episode/{episode}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// Create Opensubtitles client
		c, err := osdb.NewClient()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		c.UserAgent = OSUserAgent

		// Anonymous Login with UserAgent string will set c.Token when successful
		if err = c.LogIn("", "", ""); err != nil {
			http.NotFound(w, r)
			return
		}

		text := vars["text"]
		langs := strings.Split(vars["lang"], ",")
		
		// Fallback language always English
		fallbackLang := false;
		for _, l := range langs {
			if l == "eng" {
				fallbackLang = true
				break
			}
		}

		if fallbackLang == false {
			langs = append(langs, "eng")
		}

		log.Println("Search subtitle by text...")
		
		season, err := strconv.ParseInt(vars["season"], 10, 64)
		episode, err :=	strconv.ParseInt(vars["episode"], 10, 64)

		params := []interface{}{}
		if season == 0 && episode == 0 {
			params = []interface{}{
				c.Token,
				[]struct {
					Query  string `xmlrpc:"query"`
					Langs string `xmlrpc:"sublanguageid"`
				}{{
					text,
					strings.Join(langs, ","),
				}},
			}
		} else {
			params = []interface{}{
				c.Token,
				[]struct {
					Query  string `xmlrpc:"query"`
					Langs string `xmlrpc:"sublanguageid"`
					Season int64  `xmlrpc:"season"`
					Episode int64  `xmlrpc:"episode"`
				}{{
					text,
					strings.Join(langs, ","),
					season,
					episode,
				}},
			}
		}

		res, err := c.SearchSubtitles(&params)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		log.Println("Subtitle found.")
		io.WriteString(w, subtitleFilesList(r.Host, res, langs[0]))
	})

	routerAPI.HandleFunc(urlAPI+"subtitlesbyfile/{hash}/{base64path}/lang/{lang}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if d, err := base64.StdEncoding.DecodeString(vars["base64path"]); err == nil {
			if t, ok := torrents[vars["hash"]]; ok {

				idx := getFileByPath(string(d), t.torrent.Files())
				file := t.torrent.Files()[idx]

				path := file.DisplayPath()
				log.Println("Calculate Opensubtitles hash...")
				
				incFileClients(path, t)

				fileHash := calculateOpensubtitlesHash(file)
				log.Println("Opensubtitles hash calculated:", fileHash)

				//stop downloading the file when no connections left
				if decFileClients(path, t) <= 0 {
					stopDownloadFile(file)				
				}

				// Create Opensubtitles client
				c, err := osdb.NewClient()
				if err != nil {
					http.NotFound(w, r)
					return
				}

				c.UserAgent = OSUserAgent

				// Anonymous Login with UserAgent string will set c.Token when successful
				if err = c.LogIn("", "", ""); err != nil {
					http.NotFound(w, r)
					return
				}

				langs := strings.Split(vars["lang"], ",")
				
				// Fallback language always English
				fallbackLang := false;
				for _, l := range langs {
					if l == "eng" {
						fallbackLang = true
						break
					}
				}

				if fallbackLang == false {
					langs = append(langs, "eng")
				}

				params := []interface{}{
					c.Token,
					[]struct {
						Hash  string `xmlrpc:"moviehash"`
						Size  int64  `xmlrpc:"moviebytesize"`
						Langs string `xmlrpc:"sublanguageid"`
					}{{
						fileHash,
						file.Length(),
						strings.Join(langs, ","),
					}},
				}

				res, err := c.SearchSubtitles(&params)
				if err != nil {
					http.NotFound(w, r)
					return
				}

				if len(res) == 0 {
					http.NotFound(w, r)
					return
				}

				log.Println("Subtitle found.")
				io.WriteString(w, subtitleFilesList(r.Host, res, langs[0]))
			} else {
				log.Println("Unknown torrent:", vars["hash"])
				http.NotFound(w, r)
				return
			}
		} else {
			log.Println(err)
			http.NotFound(w, r)
			return
		}
	})

	routerAPI.HandleFunc(urlAPI+"getsubtitle/{base64path}/encode/{encode}/subtitle.srt", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if subtitleurl, err := base64.StdEncoding.DecodeString(vars["base64path"]); err == nil {
			
			zipContent, err := fetchZip(string(subtitleurl))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			
			for _, f := range zipContent.File {
				if strings.HasSuffix(strings.ToLower(f.Name), ".srt") == true {
					fileHandler, err := f.Open()
					if err != nil {
						http.NotFound(w, r)
						return
					}
					data, err := ioutil.ReadAll(fileHandler)
					if err != nil {
						http.NotFound(w, r)
						return
					}
					fileHandler.Close()

					contentType := http.DetectContentType(data)

					w.Header().Set("Content-Disposition", "filename=subtitle.srt")
					w.Header().Set("Content-Type", contentType)

					if data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
						trimmedData := bytes.Trim(data, "\xef\xbb\xbf")
						/*writeErr := ioutil.WriteFile("tmp/subtitle.srt", data, 0644)
						if writeErr != nil {
							log.Println("Subtitle save error")
						}*/
						io.WriteString(w, strings.Replace(string(trimmedData), "{\\an8}", "", -1))
					} else {
						/*writeErr := ioutil.WriteFile("tmp/subtitle.srt", []byte("\xef\xbb\xbf" + decodeData(data, vars["encode"])), 0644)
						if writeErr != nil {
							log.Println("Subtitle save error")
						}*/
						io.WriteString(w, strings.Replace(decodeData(data, vars["encode"]), "{\\an8}", "", -1))
					}
					break
				}
			}
		} else {
			//log.Println(err)
			http.NotFound(w, r)
			return
		}
	})

	routerAPI.HandleFunc(urlAPI+"torrents", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, showAllTorrent())
	})

	routerAPI.HandleFunc(urlAPI+"getmoviemagnet/imdb/{imdb}/providers/{providers}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Printf("Getting movie magnet link by this imdb id: %v\n", vars["imdb"])

		output := providers.GetMovieMagnet(vars["imdb"], "", strings.Split(vars["providers"], ","))
		if output == "" {
			log.Printf("Not found any magnet link.\n")
		} else {
			log.Printf("Magnet link found.\n")
		}

		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"getmoviemagnet/query/{query}/providers/{providers}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Printf("Getting movie magnet link by this query: %v\n", vars["query"])

		output := providers.GetMovieMagnet("", vars["query"], strings.Split(vars["providers"], ","))
		if output == "" {
			log.Printf("Not found any magnet link.\n")
		} else {
			log.Printf("Magnet link found.\n")
		}

		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"getmoviemagnet/imdb/{imdb}/query/{query}/providers/{providers}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Printf("Getting movie magnet link by this imdb id: %v, query: %v\n", vars["imdb"], vars["query"])

		output := providers.GetMovieMagnet(vars["imdb"], vars["query"], strings.Split(vars["providers"], ","))
		if output == "" {
			log.Printf("Not found any magnet link.\n")
		} else {
			log.Printf("Magnet link found.\n")
		}

		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"getshowmagnet/imdb/{imdb}/season/{season}/episode/{episode}/providers/{providers}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Printf("Getting tv show magnet link by this imdb id: %v, season: %v, episode: %v\n", vars["imdb"], vars["season"], vars["episode"])

		output := providers.GetShowMagnet(vars["imdb"], "", vars["season"], vars["episode"], strings.Split(vars["providers"], ","))
		if output == "" {
			log.Printf("Not found any magnet link.\n")
		} else {
			log.Printf("Magnet link found.\n")
		}

		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"getshowmagnet/query/{query}/season/{season}/episode/{episode}/providers/{providers}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Printf("Getting tv show magnet link by this query: %v, season: %v, episode: %v\n", vars["query"], vars["season"], vars["episode"])

		output := providers.GetShowMagnet("", vars["query"], vars["season"], vars["episode"], strings.Split(vars["providers"], ","))
		if output == "" {
			log.Printf("Not found any magnet link.\n")
		} else {
			log.Printf("Magnet link found.\n")
		}

		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"getshowmagnet/imdb/{imdb}/query/{query}/season/{season}/episode/{episode}/providers/{providers}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Printf("Getting tv show magnet link by this imdb id: %v, query: %v, season: %v, episode: %v\n", vars["imdb"], vars["query"], vars["season"], vars["episode"])

		output := providers.GetShowMagnet(vars["imdb"], vars["query"], vars["season"], vars["episode"], strings.Split(vars["providers"], ","))
		if output == "" {
			log.Printf("Not found any magnet link.\n")
		} else {
			log.Printf("Magnet link found.\n")
		}

		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"tmdbdiscover/type/{type}/genretype/{genretype}/sort/{sort}/date/{date}/lang/{lang}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("Get TMDB list by genre")

		output := providers.MirrorTmdbDiscover(vars["type"], vars["genretype"], vars["sort"], vars["date"], vars["lang"], vars["page"])
		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"tmdbsearch/type/{type}/lang/{lang}/page/{page}/text/{text}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("Get TMDB search")

		output := providers.MirrorTmdbSearch(vars["type"], vars["lang"], vars["page"], vars["text"])
		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"tmdbinfo/type/{type}/tmdbid/{tmdbid}/lang/{lang}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("Get TMDB info")

		output := providers.MirrorTmdbInfo(vars["type"], vars["tmdbid"], vars["lang"])
		io.WriteString(w, output)
	})

	routerAPI.HandleFunc(urlAPI+"receivemagnet/{todo}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		io.WriteString(w, checkReceivedMagnetHash(vars["todo"]))
	})

	routerAPI.HandleFunc(urlAPI+"websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, _ := upgrader.Upgrade(w, r, nil) // Error ignored

		for {
			// Read message from ws
            messageType, message, err := ws.ReadMessage()
            if err != nil {
                return
            }
            
            if messageType == 1 {
            	if string(message) == "stop" {
            		if err = ws.WriteMessage(1, []byte("{\"function\":\"stopserver\",\"data\": \"ok\"}")); err != nil {
	            		return
	            	}

            		go func() {
						time.Sleep(1 * time.Nanosecond)
						procQuit <- true
					}()
            	} else {
	            	value := setReceivedMagnetHash(string(message))
	            	if err = ws.WriteMessage(1, []byte("{\"function\":\"sendmagnet\",\"data\":\"" + value + "\"}")); err != nil {
	            		return
	            	}
	            }
            } else if messageType == 2 {
            	metaData, error := metainfo.Load(bytes.NewReader(message))
				if error == nil {
					spec := torrent.TorrentSpecFromMetaInfo(metaData)
					log.Println("Torrent file received:", spec.InfoHash.String())

					value := setReceivedMagnetHash(spec.InfoHash.String())
					if err = ws.WriteMessage(1, []byte("{\"function\":\"sendfile\",\"data\":\"" + value + "\"}")); err != nil {
	            		return
	            	}
				}
            }
            
        }
	})

	// Enable CORS for api urls if required
	if cors == false {
		http.Handle(urlAPI, routerAPI)
	} else {
		http.Handle(urlAPI, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(routerAPI))
	}

	// Create torrent magnet send page from main page
	sendMagnetPage := mux.NewRouter()
	sendMagnetPage.SkipClean(true)

	sendMagnetPage.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {		
		io.WriteString(w, createServerPage())
	})

	// Enable CORS for torrent sender page if required
	if cors == false {
		http.Handle("/", sendMagnetPage)
	} else {
		http.Handle("/", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(sendMagnetPage))
	}
}

func startHTTPServer(addr string, cors bool) *http.Server {
	srv := &http.Server{
		Addr: addr,
		ReadTimeout:  38 * time.Second,
		WriteTimeout: 38 * time.Second,
	}

	// Must appear
	fmt.Printf("White Raven Server Version %s Started On Address: http://127.0.0.1%s\n", version, srv.Addr)

	handleAPI(cors)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			if err == http.ErrServerClosed {
				fmt.Printf("HTTP Server Closed")
			} else {
				fmt.Printf("HTTP Server Error: %s", err)
			}
			time.Sleep(1 * time.Nanosecond)
			procQuit <- true
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}
