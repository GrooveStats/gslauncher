# Filesystem IPC

The file system layout is like this:

- requests/: request files. The filename format is <id>.json.
- responses/: response files. The filename matches the request filename.

Stale requests and response files (older than 1 minute) are removed regularly.
This assumes that SM doesn't wait for responses for more than a minute.


## Requests

### Ping

```jsonc
{
    "action": "ping",
    "payload": "something"
}
```


### Get Scores

```jsonc
{
    "action": "get-scores",
    "api-key": "topsecret",
    "hash": "somehash"
}
```


### Submit Score

```jsonc
{
    "action": "submit-score",
    "api-key": "topsecret",
    "hash": "somehash",
    "rate": 199,        // music rate x100
    "score": 9900       // score x100
}
```


## Responses

The response for ping looks like this:

```jsonc
{"payload": "something"}
```

Responses for network requests look like this:

```jsonc
{
    success: true,
    data: {}    // data returned by the endpoint
}
```


## GrooveStats Faking

The launcher can be built with the "fake" build tag to replace requests to
GrooveStats with fake responses.

```sh
go build -tags fake ./cmd/gslauncher/
```

### Get Scores

The launcher randomly returns either:
- A network error
- A response without rgpData
- A response with rpgData


### Submit Score

The response depends on the music rate:
- rate <= 33: A network error
- rate <= 66: Score added (no rpgData)
- rate <= 100: Score added (rpgData)
- rate <= 133: Score improved
- rate <= 166: Score improved
- rate > 166: Song not ranked
