package groovestats

import (
	"embed"
	"encoding/json"
	"errors"
	"time"

	"github.com/GrooveStats/gslauncher/internal/fsipc"
	"github.com/GrooveStats/gslauncher/internal/settings"
)

//go:embed fake/*.json
var fs embed.FS

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fakeNewSession() (*NewSessionResponse, error) {
	delay := time.Duration(settings.Get().FakeGsNetworkDelay)
	time.Sleep(delay * time.Second)

	if settings.Get().FakeGsNetworkError {
		return nil, errors.New("network error")
	}

	var response NewSessionResponse
	err := loadFakeData("new-session.json", &response)
	if err != nil {
		return nil, err
	}

	if !settings.Get().FakeGsRpg {
		response.ActiveEvents = response.ActiveEvents[1:]
	}
	if !settings.Get().FakeGsItl {
		response.ActiveEvents = response.ActiveEvents[:len(response.ActiveEvents)-1]
	}

	response.ServicesResult = settings.Get().FakeGsNewSessionResult
	if response.ServicesResult != "OK" {
		response.ServicesAllowed.ScoreSubmit = false
		response.ServicesAllowed.PlayerScores = false
		response.ServicesAllowed.PlayerLeaderboards = false
	}

	return &response, nil
}

func fakePlayerScores(request *fsipc.GsPlayerScoresRequest) (*PlayerScoresResponse, error) {
	delay := time.Duration(settings.Get().FakeGsNetworkDelay)
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
	delay := time.Duration(settings.Get().FakeGsNetworkDelay)
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
		p1 := response.Player1
		p2 := response.Player2

		p1.GsLeaderboard = p1.GsLeaderboard[:min(n, len(p1.GsLeaderboard))]
		p1.ExLeaderboard = p1.ExLeaderboard[:min(n, len(p1.ExLeaderboard))]
		p1.Rpg.RpgLeaderboard = p1.Rpg.RpgLeaderboard[:min(n, len(p1.Rpg.RpgLeaderboard))]
		p2.GsLeaderboard = p2.GsLeaderboard[:min(n, len(p2.GsLeaderboard))]
		p2.ExLeaderboard = p2.ExLeaderboard[:min(n, len(p2.ExLeaderboard))]
		p2.Rpg.RpgLeaderboard = p2.Rpg.RpgLeaderboard[:min(n, len(p2.Rpg.RpgLeaderboard))]
	}

	if !settings.Get().FakeGsRpg {
		response.Player1.Rpg = nil
		response.Player2.Rpg = nil
	}
	if !settings.Get().FakeGsItl {
		response.Player1.Itl = nil
		response.Player2.Itl = nil
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
	delay := time.Duration(settings.Get().FakeGsNetworkDelay)
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
		p1 := response.Player1
		p2 := response.Player2

		p1GsLeaderboard := (*p1.GsLeaderboard)[:min(n, len(*p1.GsLeaderboard))]
		p1ExLeaderboard := (*p1.ExLeaderboard)[:min(n, len(*p1.ExLeaderboard))]
		p1RpgLeaderboard := (*p1.Rpg.RpgLeaderboard)[:min(n, len(*p1.Rpg.RpgLeaderboard))]
		p2GsLeaderboard := (*p2.GsLeaderboard)[:min(n, len(*p2.GsLeaderboard))]
		p2ExLeaderboard := (*p2.ExLeaderboard)[:min(n, len(*p2.ExLeaderboard))]
		p2RpgLeaderboard := (*p2.Rpg.RpgLeaderboard)[:min(n, len(*p2.Rpg.RpgLeaderboard))]

		p1.GsLeaderboard = &p1GsLeaderboard
		p1.ExLeaderboard = &p1ExLeaderboard
		p1.Rpg.RpgLeaderboard = &p1RpgLeaderboard
		p2.GsLeaderboard = &p2GsLeaderboard
		p2.ExLeaderboard = &p2ExLeaderboard
		p2.Rpg.RpgLeaderboard = &p2RpgLeaderboard
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
	case "chart-not-ranked":
		response.Player1.IsRanked = false
		response.Player1.ScoreDelta = nil
		response.Player1.GsLeaderboard = nil
		response.Player1.ExLeaderboard = nil
		response.Player1.Rpg = nil
		response.Player1.IsRanked = false
		response.Player2.ScoreDelta = nil
		response.Player2.GsLeaderboard = nil
		response.Player2.ExLeaderboard = nil
		response.Player2.Rpg = nil
	default:
		panic("unknown submit result")
	}

	if !settings.Get().FakeGsRpg {
		response.Player1.Rpg = nil
		response.Player2.Rpg = nil
	}
	if !settings.Get().FakeGsItl {
		response.Player1.Itl = nil
		response.Player2.Itl = nil
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
