package groovestats

import (
	"embed"
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/archiveflax/gslauncher/internal/fsipc"
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

func fakePlayerScores(request *fsipc.GsPlayerScoresRequest) (*PlayerScoresResponse, error) {
	switch rand.Intn(2) {
	case 0:
		return nil, errors.New("network error")
	}

	var response PlayerScoresResponse
	err := loadFakeData("player-scores.json", &response)
	if err != nil {
		return nil, err
	}

	if request.ApiKeyPlayer1 == nil {
		response.Player1 = nil
	}

	if request.ApiKeyPlayer2 == nil {
		response.Player2 = nil
	}

	return &response, nil
}

func fakePlayerLeaderboards(request *fsipc.GsPlayerLeaderboardsRequest) (*PlayerLeaderboardsResponse, error) {
	switch rand.Intn(2) {
	case 0:
		return nil, errors.New("network error")
	}

	var response PlayerLeaderboardsResponse
	err := loadFakeData("player-leaderboards.json", &response)
	if err != nil {
		return nil, err
	}

	if request.MaxLeaderboardResults != nil {
		n := *request.MaxLeaderboardResults

		response.Player1.GsLeaderboard = response.Player1.GsLeaderboard[:n]
		response.Player1.Rpg.RpgLeaderboard = response.Player1.Rpg.RpgLeaderboard[:n]
		response.Player2.GsLeaderboard = response.Player2.GsLeaderboard[:n]
		response.Player2.Rpg.RpgLeaderboard = response.Player2.Rpg.RpgLeaderboard[:n]
	}

	if request.ApiKeyPlayer1 == nil {
		response.Player1 = nil
	}

	if request.ApiKeyPlayer2 == nil {
		response.Player2 = nil
	}

	return &response, nil
}

func fakeScoreSubmit(request *fsipc.GsScoreSubmitRequest) (*ScoreSubmitResponse, error) {
	switch rand.Intn(2) {
	case 0:
		return nil, errors.New("network error")
	}

	var response ScoreSubmitResponse
	err := loadFakeData("score-submit.json", &response)
	if err != nil {
		return nil, err
	}

	if request.MaxLeaderboardResults != nil {
		n := *request.MaxLeaderboardResults

		p1GsLeaderboard := (*response.Player1.GsLeaderboard)[:n]
		p1RpgLeaderboard := (*response.Player1.Rpg.RpgLeaderboard)[:n]
		p2GsLeaderboard := (*response.Player2.GsLeaderboard)[:n]
		p2RpgLeaderboard := (*response.Player2.Rpg.RpgLeaderboard)[:n]

		response.Player1.GsLeaderboard = &p1GsLeaderboard
		response.Player1.Rpg.RpgLeaderboard = &p1RpgLeaderboard
		response.Player2.GsLeaderboard = &p2GsLeaderboard
		response.Player2.Rpg.RpgLeaderboard = &p2RpgLeaderboard
	}

	if request.Player1 == nil {
		response.Player1 = nil
	} else {
		switch {
		case request.Player1.Rate <= 33:
			response.Player1.Result = "score-added"
			response.Player1.ScoreDelta = nil
			response.Player1.Rpg = nil
		case request.Player1.Rate <= 66:
			response.Player1.Result = "score-added"
			response.Player1.ScoreDelta = nil
			response.Player1.Rpg.Result = "score-added"
			response.Player1.Rpg.ScoreDelta = nil
			response.Player1.Rpg.RateDelta = nil
		case request.Player1.Rate <= 100:
			response.Player1.Result = "score-improved"
			response.Player1.Rpg.Result = "score-improved"
		case request.Player1.Rate <= 133:
			zero := 0

			response.Player1.Result = "score-not-improved"
			response.Player1.ScoreDelta = &zero
			response.Player1.Rpg.Result = "score-not-improved"
			response.Player1.Rpg.ScoreDelta = &zero
			response.Player1.Rpg.RateDelta = &zero
		default:
			response.Player1.Result = "song-not-ranked"
			response.Player1.ScoreDelta = nil
			response.Player1.GsLeaderboard = nil
			response.Player1.Rpg.Result = "song-not-ranked"
			response.Player1.Rpg.ScoreDelta = nil
			response.Player1.Rpg.RateDelta = nil
			response.Player1.Rpg.Progress = nil
			response.Player1.Rpg.RpgLeaderboard = nil
		}
	}

	if request.Player2 == nil {
		response.Player2 = nil
	} else {
		switch {
		case request.Player2.Rate <= 33:
			response.Player2.Result = "score-added"
			response.Player2.ScoreDelta = nil
			response.Player2.Rpg = nil
		case request.Player2.Rate <= 66:
			response.Player2.Result = "score-added"
			response.Player2.ScoreDelta = nil
			response.Player2.Rpg.Result = "score-added"
			response.Player2.Rpg.ScoreDelta = nil
			response.Player2.Rpg.RateDelta = nil
		case request.Player2.Rate <= 100:
			response.Player2.Result = "score-improved"
			response.Player2.Rpg.Result = "score-improved"
		case request.Player2.Rate <= 133:
			zero := 0

			response.Player2.Result = "score-not-improved"
			response.Player2.ScoreDelta = &zero
			response.Player2.Rpg.Result = "score-not-improved"
			response.Player2.Rpg.ScoreDelta = &zero
			response.Player2.Rpg.RateDelta = &zero
		default:
			response.Player2.Result = "song-not-ranked"
			response.Player2.ScoreDelta = nil
			response.Player2.GsLeaderboard = nil
			response.Player2.Rpg.Result = "song-not-ranked"
			response.Player2.Rpg.ScoreDelta = nil
			response.Player2.Rpg.RateDelta = nil
			response.Player2.Rpg.Progress = nil
			response.Player2.Rpg.RpgLeaderboard = nil
		}
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
