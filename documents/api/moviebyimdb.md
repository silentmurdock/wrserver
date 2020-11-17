# Get Movie Magnet By IMDB Id

Search magnet links for movie by IMDB id.

**URL** : `/api/getmoviemagnet/imdb/{imdb}/providers/{providers}`

**Method** : `GET`

**Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `imdb` | string | Internet Movie Database identifier.|
| `providers` | string | Public torrent site identifiers separated by a comma.|

**Supported Providers** :

| Provider | Type | Website |
| --- | --- | --- |
| `pt` | string | POPCORN TIME|
| `yts` | string | YTS|
| `itorrent` | string | ITORRENT|
| `rarbg` | string | RARBG|
| `pto` | string | POPCORN TIME ONLINE|

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

`GET http://localhost:9000/api/getmoviemagnet/imdb/tt7556122/providers/pt,yts,itorrent`

**Success Response** :

```json
{
    "success": true,
    "results": [{
        "hash": "E9C667FBED735AD0C6602748EB5CDC7709CC42AA",
        "quality": "1080p",
        "size": "2491081032",
        "provider": "YTS",
        "lang": "en",
        "title": "The Old Guard",
        "seeds": "1033",
        "peers": "254"
    }, {
        "hash": "6ECBD1C6AD35857CA0ED802D7DE41C30A90920C5",
        "quality": "720p",
        "size": "1213328261",
        "provider": "YTS",
        "lang": "en",
        "title": "The Old Guard",
        "seeds": "993",
        "peers": "184"
    }, {
        "hash": "eddc84cb5532ee016ddb5192ba08e0b88e05a573",
        "quality": "WEBRIP",
        "size": "2040109466",
        "provider": "ITORRENT",
        "lang": "hu",
        "title": "The.Old.Guard.2020.NF.WEBRip.x264.HuN-prldm",
        "seeds": "20",
        "peers": "11"
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