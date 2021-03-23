package groovestats

import (
	"embed"
	"encoding/json"
	"errors"
	"time"

	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/settings"
)

//go:embed fake/*.json
var fs embed.FS

func fakeNewSession() (*NewSessionResponse, error) {
	delay := time.Duration(settings.Get().FakeGsNetDelay)
	time.Sleep(delay * time.Second)

	if settings.Get().FakeGsNetworkError {
		return nil, errors.New("network error")
	}

	var filename string

	if settings.Get().FakeGsDdos {
		filename = "new-session-ddos.json"
	} else {
		filename = "new-session.json"
	}

	var response NewSessionResponse
	err := loadFakeData(filename, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func fakePlayerScores(request *fsipc.GsPlayerScoresRequest) (*PlayerScoresResponse, error) {
	delay := time.Duration(settings.Get().FakeGsNetDelay)
	time.Sleep(delay * time.Second)

	if settings.Get().FakeGsNetworkError {
		return nil, errors.New("network error")
	}

	var response PlayerScoresResponse
	err := loadFakeData("player-scores.json", &response)
	if err != nil {
		return nil, err
	}

	if request.Player1 == nil {
		response.Player1 = nil
	} else {
		response.Player1.ChartHash = request.Player1.ChartHash
	}

	if request.Player2 == nil {
		response.Player2 = nil
	} else {
		response.Player2.ChartHash = request.Player2.ChartHash
	}

	return &response, nil
}

func fakePlayerLeaderboards(request *fsipc.GsPlayerLeaderboardsRequest) (*PlayerLeaderboardsResponse, error) {
	delay := time.Duration(settings.Get().FakeGsNetDelay)
	time.Sleep(delay * time.Second)

	if settings.Get().FakeGsNetworkError {
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

	if !settings.Get().FakeGsRpg {
		response.Player1.Rpg = nil
		response.Player2.Rpg = nil
	}

	if request.Player1 == nil {
		response.Player1 = nil
	} else {
		response.Player1.ChartHash = request.Player1.ChartHash
	}

	if request.Player2 == nil {
		response.Player2 = nil
	} else {
		response.Player2.ChartHash = request.Player2.ChartHash
	}

	return &response, nil
}

func fakeScoreSubmit(request *fsipc.GsScoreSubmitRequest) (*ScoreSubmitResponse, error) {
	delay := time.Duration(settings.Get().FakeGsNetDelay)
	time.Sleep(delay * time.Second)

	if settings.Get().FakeGsNetworkError {
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

	switch settings.Get().FakeGsSubmitResult {
	case "score-added":
		response.Player1.Result = "score-added"
		response.Player1.ScoreDelta = nil
		response.Player1.Rpg.Result = "score-added"
		response.Player1.Rpg.ScoreDelta = nil
		response.Player1.Rpg.RateDelta = nil
		response.Player2.Result = "score-added"
		response.Player2.ScoreDelta = nil
		response.Player2.Rpg.Result = "score-added"
		response.Player2.Rpg.ScoreDelta = nil
		response.Player2.Rpg.RateDelta = nil
	case "improved":
		response.Player1.Result = "improved"
		response.Player1.Rpg.Result = "improved"
		response.Player2.Result = "improved"
		response.Player2.Rpg.Result = "improved"
	case "score-not-improved":
		zero := 0

		response.Player1.Result = "score-not-improved"
		response.Player1.ScoreDelta = &zero
		response.Player1.Rpg.Result = "score-not-improved"
		response.Player1.Rpg.ScoreDelta = &zero
		response.Player1.Rpg.RateDelta = &zero
		response.Player2.Result = "score-not-improved"
		response.Player2.ScoreDelta = &zero
		response.Player2.Rpg.Result = "score-not-improved"
		response.Player2.Rpg.ScoreDelta = &zero
		response.Player2.Rpg.RateDelta = &zero
	case "score-not-ranked":
		response.Player1.Result = "song-not-ranked"
		response.Player1.ScoreDelta = nil
		response.Player1.GsLeaderboard = nil
		response.Player1.Rpg.Result = "song-not-ranked"
		response.Player1.Rpg.ScoreDelta = nil
		response.Player1.Rpg.RateDelta = nil
		response.Player1.Rpg.Progress = nil
		response.Player1.Rpg.RpgLeaderboard = nil
		response.Player2.Result = "song-not-ranked"
		response.Player2.ScoreDelta = nil
		response.Player2.GsLeaderboard = nil
		response.Player2.Rpg.Result = "song-not-ranked"
		response.Player2.Rpg.ScoreDelta = nil
		response.Player2.Rpg.RateDelta = nil
		response.Player2.Rpg.Progress = nil
		response.Player2.Rpg.RpgLeaderboard = nil
	default:
		panic("unknown submit result")
	}

	if !settings.Get().FakeGsRpg {
		response.Player1.Rpg = nil
		response.Player2.Rpg = nil
	}

	if request.Player1 == nil {
		response.Player1 = nil
	} else {
		response.Player1.ChartHash = request.Player1.ChartHash
	}

	if request.Player2 == nil {
		response.Player2 = nil
	} else {
		response.Player2.ChartHash = request.Player2.ChartHash
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
