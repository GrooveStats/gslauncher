# Filesystem IPC

The file system layout is like this:

- requests/: request files. The filename format is <id>.json.
- responses/: response files. The filename matches the request filename.

Stale requests and response files (older than 1 minute) are removed regularly.
This assumes that SM doesn't wait for responses for more than a minute.


## Requests

### Ping

```jsonc
{"action": "ping"}
```


### GrooveStats: New Session

```jsonc
{"action": "groovestats/new-session"}
```


### GrooveStats: Player Scores

```jsonc
{
    "action": "groovestats/player-scores",
    "chart": "somehash",                        // the hash of the chart
    "api-key-player-1": "topsecret",            // optional
    "api-key-player-2": "topsecret"             // optional
}
```

At least one of the two API keys has to be provided.


### GrooveStats: Player Leaderboards

```jsonc
{
    "action": "groovestats/player-leaderboards",
    "chart": "somehash",                        // the hash of the chart
    "max-leaderboard-results": 10,              // optional
    "api-key-player-1": "topsecret",            // optional
    "api-key-player-2": "topsecret"             // optional
}
```

At least one of the two API keys has to be provided.


### Submit Score

```jsonc
{
    "action": "groovestats/submit-score",
    "api-key": "topsecret",
    "hash": "somehash",
    "rate": 199,        // music rate x100
    "score": 9900       // score x100
}
```


## Responses

The response for ping looks like this:

```jsonc
{}
```

Responses for network requests look like this:

```jsonc
{
    success: true,
    data: {}    // data returned by the endpoint
}
```


## GrooveStats Faking

The debug build of the launcher adds support for the "Fake GrooveStats
Requests" setting. It replaces requests to GrooveStats with predetermined fake
responses.

```sh
go build -tags debug ./cmd/gslauncher/
```


### New Session

The launcher randomly returns either:
- A network error
- A response that indicates that all services are up
- A response that indicates that all services are down (DDoS scenario)


### Player Scores

The launcher randomly returns either:
- A network error
- A leaderboard response for the requested players


### Player Leaderboards

The launcher randomly returns either:
- A network error
- A leaderboard response for the requested players


### Submit Score

The response depends on the music rate:
- rate <= 33: A network error
- rate <= 66: Score added (no rpgData)
- rate <= 100: Score added (rpgData)
- rate <= 133: Score improved
- rate <= 166: Score improved
- rate > 166: Song not ranked
