# Get Movie Magnet By Query Text

Search magnet links for movie by query text.

**URL** : `/api/getmoviemagnet/query/{query}/providers/{providers}`

**Method** : `GET`

**Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `query` | string | Text query to search. **This value should be URI encoded.**|
| `providers` | string | Public torrent site identifiers separated by a comma.|

**Text Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `title` | string | Movie title.|
| `releaseyear` | string | Release year.|

**Supported Providers** :

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
| `size` | string | Torrent data size in bytes.|
| `provider` | string | Source of the magnet link.|
| `lang` | string | ISO 639-1 two-letter language code.|
| `title` | string | Movie title.|
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

`GET http://localhost:9000/api/getmoviemagnet/query/title%3DProject%20Power%26releaseyear%3D2020/providers/1337x`

**Success Response** :

```json
{
    "success": true,
    "results": [{
        "hash": "937712B36BEB9F9E4A878C5ACCE72250C646443D",
        "quality": "720p",
        "size": "1073741824",
        "provider": "1337X",
        "lang": "en",
        "title": "Project Power 2020 720p WEBRip YTS YIFY",
        "seeds": "4659",
        "peers": "1976"
    }, {
        "hash": "84B29629AFDA75ACD6C1AE90F45D8E822954B339",
        "quality": "1080p",
        "size": "2254857830",
        "provider": "1337X",
        "lang": "en",
        "title": "Project Power 2020 1080p WEBRip 5.1 YTS YIFY",
        "seeds": "4016",
        "peers": "2404"
    }, {
        "hash": "4758D81B183836C06FCEB1601C4688B2E6AC6BD5",
        "quality": "720p",
        "size": "835085926",
        "provider": "1337X",
        "lang": "en",
        "title": "Project.Power.2020.720p.NF.WEBrip.800MB.x264-GalaxyRG",
        "seeds": "1756",
        "peers": "440"
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