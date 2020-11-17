# White Raven Server

White Raven Server is a REST-like API controlled torrent client application to find movies and tv shows from various sources and stream them over http connection. Mainly created for [White Raven](https://github.com/silentmurdock/whiteraven), which is a torrent player application for  Samsung Smart TV E, F, H series.

## HTTP API Functions
### Server and Client related
* [Get server information](documents/api/about.md) : `GET /api/about`
* [Stop server](documents/api/stop.md) : `GET /api/stop`
* [Restart torrent client](documents/api/restart.md) : `GET /api/restart/downrate/{downrate}/uprate/{uprate}`

### Torrent related
* [Add torrent by hash](documents/api/add.md) : `GET /api/add/{hash}`
* [Delete torrent by hash](documents/api/delete.md) : `GET /api/delete/{hash}`
* [Delete all running torrents](documents/api/deleteall.md) : `GET /api/deleteall`
* [Get all running torrents](documents/api/torrents.md) : `GET /api/torrents`
* [Get running torrent statistics by hash](documents/api/stats.md) : `GET /api/stats/{hash}`
* [Stream or download the selected file](documents/api/get.md) : `GET /api/get/{hash}/{base64path}`

### Movie or TV Show related
* [Get movie magnet by IMDB id](documents/api/moviebyimdb.md) : `GET /api/getmoviemagnet/imdb/{imdb}/providers/{providers}`
* [Get movie magnet by query text](documents/api/moviebytext.md) : `GET /api/getmoviemagnet/query/{query}/providers/{providers}`
* [Get movie magnet by IMDB id and query text at once](documents/api/moviebyboth.md) : `GET /api/getmoviemagnet/imdb/{imdb}/query/{query}/providers/{providers}`
* [Get show magnet by IMDB id](documents/api/showbyimdb.md) : `GET /api/getshowmagnet/imdb/{imdb}/season/{season}/episode/{episode}/providers/{providers}`
* [Get show magnet by query text](documents/api/showbytext.md) : `GET /api/getshowmagnet/query/{query}/season/{season}/episode/{episode}/providers/{providers}`
* [Get show magnet by IMDB id and query text at once](documents/api/showbyboth.md) : `GET /api/getshowmagnet/imdb/{imdb}/query/{query}/season/{season}/episode/{episode}/providers/{providers}`
* [Discover movies or tv shows](documents/api/tmdbdiscover.md) : `GET /api/tmdbdiscover/type/{type}/genretype/{genretype}/sort/{sort}/date/{date}/lang/{lang}/page/{page}`
* [Search movies or tv shows by query text](documents/api/tmdbsearch.md) : `GET /api/tmdbsearch/type/{type}/lang/{lang}/page/{page}/text/{text}`
* [Get more info about movie or tv show by TMDB id](documents/api/tmdbinfo.md) : `GET /api/tmdbinfo/type/{type}/tmdbid/{tmdbid}/lang/{lang}`

### Subtitle related
* [Search subtitles by IMDB id](documents/api/subtitlesbyimdb.md) : `GET /api/subtitlesbyimdb/{imdb}/lang/{lang}/season/{season}/episode/{episode}`
* [Search subtitles by query text](documents/api/subtitlesbytext.md) : `GET /api/subtitlesbytext/{text}/lang/{lang}/season/{season}/episode/{episode}`
* [Search subtitles by inner file hash](documents/api/subtitlesbyhash.md) : `GET /api/subtitlesbyfile/{hash}/{base64path}/lang/{lang}`
* [Download subtitle file](documents/api/getsubtitle.md) : `GET /api/getsubtitle/{base64path}/encode/{encode}/subtitle.srt`
<br/>

## Command-Line Arguments
* **-background** run the server in the background
* **-cors** enable CORS
* **-dir** `string` specify the directory where files will be downloaded to if storagetype is set to "piecefile" or "file"
* **-downrate** `int` download speed rate in Kbps (`default 4096`)
* **-host** `string` listening server ip
* **-log** enable log messages
* **-maxconn** `int` max connections per torrent (`default 40`)
* **-memorysize** `int` specify the storage memory size in MB if storagetype is set to "memory" (`minimum 64`) (`default 64`)
* **-nodht** disable dht
* **-osuseragent**`string` set external OpenSubtitles user agent
* **-port** `int` listening port (`default 9000`)
* **-storagetype**`string` select storage type (must be set to "memory" or "piecefile" or "file")
* **-tmdbkey**`string` set external TMDB API key
* **-uprate** `int` upload speed rate in Kbps (`default 256`)
<br/>

## Build Instructions

### Build On Windows
**Download:**
```
$ go get -v -u github.com/silentmurdock/wrserver
```
**Build in vendor mode for Samsung Smart TV E, F, H ARM series:**
```
$ set GOOS=linux
$ set GOARCH=arm
$ set GOARM=7
$ go build -ldflags="-s -w" -mod=vendor -o built\linux\arm7\wrserver
```
**Build in vendor mode for Windows i386:**
```
$ set GOOS=windows
$ set GOARCH=386
$ go build -ldflags="-s -w" -mod=vendor -o built\windows\i386\wrserver.exe
```
**Build in vendor mode for Linux i386:**
```
$ set GOOS=linux
$ set GOARCH=386
$ go build -ldflags="-s -w" -mod=vendor -o built\linux\i386\wrserver
```

### Build On Linux
**Download:**
```
$ go get -v -u github.com/silentmurdock/wrserver
```
**Build in vendor mode for Samsung Smart TV E, F, H ARM series:**
```
$ env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -mod=vendor -o built/linux/arm7/wrserver
```
**Build in vendor mode for Windows i386:**
```
$ env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -mod=vendor -o built/windows/i386/wrserver.exe
```
**Build in vendor mode for Linux x64:**
```
$ env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -mod=vendor -o built/linux/x64/wrserver
```
<br/>

## Run The Server

### Run The Server On Windows
**Serve torrent data with i386 executable from memory:**
```
$ built\windows\i386\wrserver -storagetype="memory"
```
**Serve torrent data with i386 executable from local disk:**
```
$ built\windows\i386\wrserver -storagetype="file" -dir="downloads"
```

### Run The Server On Linux
**Serve torrent data with x64 executable from memory:**
```
$ ./built/linux/x64/wrserver -storagetype="memory"
```
**Serve torrent data with x64 executable from local disk:**
```
$ ./built/linux/x64/wrserver -storagetype="file" -dir="downloads"
```
<br/>

## License
[GNU GENERAL PUBLIC LICENSE Version 3](LICENSE)