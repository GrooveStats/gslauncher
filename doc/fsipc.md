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
    "protocol": 1
}
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


### GrooveStats: Score Submit

```jsonc
{
    "action": "groovestats/score-submit",
    "chart": "somehash",
    "max-leaderboard-results": 10,              // optional
    "player1": {                                // optional
        "api-key": "topsecret",
        "profile-name": "domp",
        "rate": 100,                            // music rate x100
        "comment": "C715, Reverse, Overhead, Cel",
        "score": 10000                          // score x100
    },
    "player2": {                                // optional
        "api-key": "topsecret",
        "profile-name": "natano",
        "rate": 199,                            // music rate x100
        "comment": "C675, Overhead, Cel",
        "score": 8630                           // score x100
    }
}
```

Data for at least one player has to be provided.


## Responses

The response for ping looks like this:

```jsonc
{
    "version": {
        "major": 1,
        "minor": 0,
        "patch": 0
    }
}
```

Responses for network requests look like this:

```jsonc
{
    success: true,
    data: {}    // data returned by the endpoint
}
```


## GrooveStats Simulation

The debug build of the launcher adds support for the "Simulate GrooveStats
Requests" setting. It replaces requests to GrooveStats with predetermined fake
responses.

```sh
go build -tags debug ./cmd/gslauncher/
```
