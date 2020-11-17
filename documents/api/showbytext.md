# Get Show Magnet By Query Text

Search magnet links for tv show by query text.

**URL** : `/api/getshowmagnet/query/{query}/season/{season}/episode/{episode}/providers/{providers}`

**Method** : `GET`

**Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `query` | string | Text query to search. **This value should be URI encoded.**|
| `season` | integer | Season number.|
| `episode` | integer | Episode number.|
| `providers` | string | Public torrent site identifiers separated by a comma.|

**Text Query Parameters** :

| Parameter | Type | Description |
| --- | --- | --- |
| `title` | string | Show title.|

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

`GET http://localhost:9000/api/getshowmagnet/query/title%3Dsupernatural/season/15/episode/17/providers/1337x`

**Success Response** :

```json
{
    "success": true,
    "results": [{
        "hash": "173FD9883257250564563BBAC21DFC2312B57FD8",
        "quality": "720p",
        "season": "15",
        "episode": "17",
        "size": "754974720",
        "provider": "1337X",
        "lang": "en",
        "title": "Supernatural.S15E17.720p.HDTV.x264-SYNCOPY",
        "seeds": "565",
        "peers": "45"
    }, {
        "hash": "819FD741391C0B81ED92AE9D4F0CE9049DE90561",
        "quality": "HDTV",
        "season": "15",
        "episode": "17",
        "size": "247149363",
        "provider": "1337X",
        "lang": "en",
        "title": "Supernatural.S15E17.HDTV.x264-PHOENiXTGx",
        "seeds": "549",
        "peers": "146"
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