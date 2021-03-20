package groovestats

import (
	"embed"
	"encoding/json"
	"errors"
	"math/rand"
)

//go:embed fake/*.json
var fs embed.FS

func fakeNewSession() (*NewSessionResponse, error) {
	var filename string

	switch rand.Intn(3) {
	case 0:
		return nil, errors.New("network error")
	case 1:
		filename = "new-session.json"
	case 2:
		filename = "new-session-ddos.json"
	}

	var response NewSessionResponse
	err := loadFakeData(filename, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func fakePlayerScores(chart string, apiKeyPlayer1, apiKeyPlayer2 *string) (*PlayerScoresResponse, error) {
	switch rand.Intn(2) {
	case 0:
		return nil, errors.New("network error")
	}

	var response PlayerScoresResponse
	err := loadFakeData("player-scores.json", &response)
	if err != nil {
		return nil, err
	}

	if apiKeyPlayer1 == nil {
		response.Player1 = nil
	}

	if apiKeyPlayer2 == nil {
		response.Player2 = nil
	}

	return &response, nil
}

func fakePlayerLeaderboards(chart string, maxLeaderboardResults *int, apiKeyPlayer1, apiKeyPlayer2 *string) (*PlayerLeaderboardsResponse, error) {
	switch rand.Intn(2) {
	case 0:
		return nil, errors.New("network error")
	}

	var response PlayerLeaderboardsResponse
	err := loadFakeData("player-leaderboards.json", &response)
	if err != nil {
		return nil, err
	}

	if maxLeaderboardResults != nil {
		p1GsLeaderboard := (*response.Player1.GsLeaderboard)[:*maxLeaderboardResults]
		p1RpgLeaderboard := (*response.Player1.Rpg.RpgLeaderboard)[:*maxLeaderboardResults]
		p2GsLeaderboard := (*response.Player2.GsLeaderboard)[:*maxLeaderboardResults]
		p2RpgLeaderboard := (*response.Player2.Rpg.RpgLeaderboard)[:*maxLeaderboardResults]

		response.Player1.GsLeaderboard = &p1GsLeaderboard
		response.Player1.Rpg.RpgLeaderboard = &p1RpgLeaderboard
		response.Player2.GsLeaderboard = &p2GsLeaderboard
		response.Player2.Rpg.RpgLeaderboard = &p2RpgLeaderboard
	}

	if apiKeyPlayer1 == nil {
		response.Player1 = nil
	}

	if apiKeyPlayer2 == nil {
		response.Player2 = nil
	}

	return &response, nil
}

func fakeAutoSubmitScore(hash string, rate int, score int) (*AutoSubmitScoreResponse, error) {
	var filename string

	switch {
	case rate <= 33:
		return nil, errors.New("network error")
	case rate <= 66:
		filename = "auto-submit-score-added.json"
	case rate <= 100:
		filename = "auto-submit-score-added-rpg.json"
	case rate <= 133:
		filename = "auto-submit-score-improved.json"
	case rate <= 166:
		filename = "auto-submit-score-not-improved.json"
	default:
		filename = "auto-submit-score-not-ranked.json"
	}

	var response AutoSubmitScoreResponse
	err := loadFakeData(filename, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func loadFakeData(filename string, response interface{}) error {
	data, err := fs.ReadFile("fake/" + filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, response)
}
