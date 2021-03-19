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

func fakeGetScores(hash string) (*GetScoresResponse, error) {
	var filename string

	switch rand.Intn(3) {
	case 0:
		return nil, errors.New("network error")
	case 1:
		filename = "get-scores.json"
	case 2:
		filename = "get-scores-rpg.json"
	}

	var response GetScoresResponse
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
