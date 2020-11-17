# Get Show Magnet By IMDB Id

Search magnet links for tv show by IMDB id.

**URL** : `/api/getshowmagnet/imdb/{imdb}/season/{season}/episode/{episode}/providers/{providers}`

**Method** : `GET`

**Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `imdb` | string | Internet Movie Database identifier.|
| `season` | integer | Season number.|
| `episode` | integer | Episode number.|
| `providers` | string | Public torrent site identifiers separated by a comma.|

**Supported Providers** :

| Provider | Type | Website |
| --- | --- | --- |
| `pt` | string | POPCORN TIME|
| `eztv` | string | EZTV|
| `itorrent` | string | ITORRENT|
| `rarbg` | string | RARBG|

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

`GET http://localhost:9000/api/getshowmagnet/imdb/tt4574334/season/3/episode/7/providers/eztv`

**Success Response** :

```json
{
    "success": true,
    "results": [{
        "hash": "55cdab9646b40f7b25053ede3a8bd7dbf195deff",
        "quality": "720p",
        "season": "3",
        "episode": "7",
        "size": "1763716626",
        "provider": "EZTV",
        "lang": "",
        "title": "Stranger Things S03E07 720p WEBRip X264-METCON EZTV",
        "seeds": "21",
        "peers": "2"
    }, {
        "hash": "27bc43954550b168566a117f5ed1cceaa6301c26",
        "quality": "480p",
        "season": "3",
        "episode": "7",
        "size": "296508204",
        "provider": "EZTV",
        "lang": "",
        "title": "Stranger Things S03E07 480p x264-mSD EZTV",
        "seeds": "20",
        "peers": "2"
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