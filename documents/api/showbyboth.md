# Get Show Magnet By IMDB Id And Query Text

Search magnet links for tv show by IMDB id and query text at once.

**URL** : `/api/getshowmagnet/imdb/{imdb}/query/{query}/season/{season}/episode/{episode}/providers/{providers}`

**Method** : `GET`

**Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `imdb` | string | Internet Movie Database identifier.|
| `query` | string | Text query to search. **This value should be URI encoded.**|
| `season` | integer | Season number.|
| `episode` | integer | Episode number.|
| `providers` | string | Public torrent site identifiers separated by a comma.|

**Text Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `title` | string | Show title.|

**Supported Providers For IMDB id** :

| Provider | Type | Website |
| --- | --- | --- |
| `pt` | string | POPCORN TIME|
| `eztv` | string | EZTV|
| `itorrent` | string | ITORRENT|
| `rarbg` | string | RARBG|

**Supported Providers For Query Text** :

| Provider | Type | Website |
| --- | --- | --- |
| `1337x` | string | 1337X|

## Success Response

**Code** : `200 OK`

**Main Object** :

| Name | Type | Description |
| --- | --- | --- |
| `success` | bool | Indicates whether the query was successful.|
| `results` | array[object] | Array of objects.|

**Object [ results ]** :

| Name | Type | Description |
| --- | --- | --- |
| `hash` | string | 40 characters long infohash.|
| `quality` | string | Video quality.|
| `season` | string | Season number.|
| `episode` | string | Episode number.|
| `size` | string | Torrent data size in bytes.|
| `provider` | string | Source of the magnet link.|
| `lang` | string | ISO 639-1 two-letter language code.|
| `title` | string | Show title.|
| `seeds` | string | Currently available seeds.|
| `peers` | string | Currently available peers.|

## Error Response

**Code** : `404 (Not Found)`

**Main Object** :

| Name | Type | Description |
| --- | --- | --- |
| `success` | bool | Indicates whether the query was successful.|
| `message` | string | Text message that describes the response.|

## Examples

**Request** :

```
GET http://localhost:9000/api/getshowmagnet/imdb/tt3743822/query/title%3DFear%20the%20Walking%20Dead/season/6/episode/4/providers/eztv,1337x
```

**Success Response** :

```json
{
    "success": true,
    "results": [{
        "hash": "6CAADA94534664385C82F6E75699817E464D89FC",
        "quality": "720p",
        "season": "6",
        "episode": "4",
        "size": "993420902",
        "provider": "1337X",
        "lang": "en",
        "title": "Fear.the.Walking.Dead.S06E04.720p.WEB.H264-CAKES",
        "seeds": "400",
        "peers": "38"
    }, {
        "hash": "218161f75a8e1cd6a6f2e8184901f070011ca853",
        "quality": "720p",
        "season": "6",
        "episode": "4",
        "size": "254628381",
        "provider": "EZTV",
        "lang": "",
        "title": "Fear the Walking Dead S06E04 720p HEVC x265-MeGusta EZTV",
        "seeds": "195",
        "peers": "23"
    }]
}
```

**Error Response** :

```json
{
    "success": false,
    "message": "No magnet links found."
}
```