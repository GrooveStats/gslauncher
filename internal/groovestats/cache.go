package groovestats

import (
	"time"

	"github.com/GrooveStats/gslauncher/internal/fsipc"
)

type playerScoresCacheEntry struct {
	data      *playerScoresPlayerData
	retrieved time.Time
}

const timeout = 15 * time.Minute

func (client *Client) addPlayerScores(request *fsipc.GsPlayerScoresRequest, response *PlayerScoresResponse) {
	now := time.Now()

	setPlayerData := func(chartHash, apiKey string, data *playerScoresPlayerData) {
		key := "player-scores:" + chartHash + "," + apiKey

		client.cache.Add(key, playerScoresCacheEntry{
			data:      data,
			retrieved: now,
		})
	}

	if request.Player1 != nil {
		player := request.Player1
		setPlayerData(player.ChartHash, player.ApiKey, response.Player1)
	}

	if request.Player2 != nil {
		player := request.Player2
		setPlayerData(player.ChartHash, player.ApiKey, response.Player2)
	}
}

func (client *Client) getPlayerScores(request *fsipc.GsPlayerScoresRequest) *PlayerScoresResponse {
	now := time.Now()

	getPlayerData := func(chartHash, apiKey string) *playerScoresPlayerData {
		key := "player-scores:" + chartHash + "," + apiKey

		iEntry, ok := client.cache.Get(key)
		if !ok {
			return nil
		}

		entry := iEntry.(playerScoresCacheEntry)
		if now.Sub(entry.retrieved) > timeout {
			client.cache.Remove(key)
			return nil
		}

		return entry.data
	}

	response := PlayerScoresResponse{Cached: true}

	if request.Player1 != nil {
		player := request.Player1
		data := getPlayerData(player.ChartHash, player.ApiKey)

		if data == nil {
			return nil
		}

		response.Player1 = data
	}

	if request.Player2 != nil {
		player := request.Player2
		data := getPlayerData(player.ChartHash, player.ApiKey)

		if data == nil {
			return nil
		}

		response.Player2 = data
	}

	return &response
}

func (client *Client) removePlayerScores(request *fsipc.GsScoreSubmitRequest) {
	removePlayerData := func(chartHash, apiKey string) {
		key := "player-scores:" + chartHash + "," + apiKey
		client.cache.Remove(key)
	}

	if request.Player1 != nil {
		player := request.Player1
		removePlayerData(player.ChartHash, player.ApiKey)
	}

	if request.Player2 != nil {
		player := request.Player2
		removePlayerData(player.ChartHash, player.ApiKey)
	}
}
