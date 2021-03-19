package groovestats

import (
	"embed"
	"encoding/json"
	"errors"
	"math/rand"
)

//go:embed fake/*.json
var fs embed.FS

func fakeAutoSubmitScore(hash string, rate int, score int) (*AutoSubmitScoreResponse, error) {
	var filename string

	switch {
	case rate <= 33:
		return nil, errors.New("network error")
	case rate <= 66:
		filename = "fake/auto-submit-score-added.json"
	case rate <= 100:
		filename = "fake/auto-submit-score-added-rpg.json"
	case rate <= 133:
		filename = "fake/auto-submit-score-improved.json"
	case rate <= 166:
		filename = "fake/auto-submit-score-not-improved.json"
	default:
		filename = "fake/auto-submit-score-not-ranked.json"
	}

	data, err := fs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var response AutoSubmitScoreResponse
	err = json.Unmarshal(data, &response)
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
		filename = "fake/get-scores.json"
	case 2:
		filename = "fake/get-scores-rpg.json"
	}

	data, err := fs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var response GetScoresResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
