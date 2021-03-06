package fsipc

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func writeRequest(filename string, data []byte) error {
	tmpfile := filename + ".new"

	err := os.WriteFile(tmpfile, data, 0700)
	if err != nil {
		return err
	}

	return os.Rename(tmpfile, filename)
}

func TestFsipc(t *testing.T) {
	dir := t.TempDir()

	ipc, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	getRequest := func(t *testing.T) interface{} {
		var request interface{}
		var ok bool

		select {
		case request, ok = <-ipc.GsPlayerScoresRequests:
			if !ok {
				t.Fatal("player-scores requests channel closed prematurely")
			}
		case request, ok = <-ipc.GsPlayerLeaderboardsRequests:
			if !ok {
				t.Fatal("player-leaderboards requests channel closed prematurely")
			}
		case request, ok = <-ipc.Requests:
			if !ok {
				t.Fatal("requests channel closed prematurely")
			}
		case <-time.After(10 * time.Second):
			t.Fatal("request not received")
		}

		return request
	}

	t.Run("PingRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "bda6a8a9d7924c149697e13b93aa68bf.json")
			err := writeRequest(filename, []byte(`{
				"action": "ping",
				"protocol": 1
			}`))
			if err != nil {
				t.Fatal(err)
			}
		}()

		request := getRequest(t)
		pingRequest, ok := request.(*PingRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := PingRequest{
			Id:       "bda6a8a9d7924c149697e13b93aa68bf",
			Protocol: 1,
		}
		if !reflect.DeepEqual(*pingRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GsNewSessionRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "b787ae38ea0c465e8d853015db940915.json")
			err := writeRequest(filename, []byte(`{
				"action": "groovestats/new-session",
				"chartHashVersion": 3
			}`))
			if err != nil {
				t.Fatal(err)
			}
		}()

		request := getRequest(t)
		newSessionRequest, ok := request.(*GsNewSessionRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsNewSessionRequest{
			Id:               "b787ae38ea0c465e8d853015db940915",
			ChartHashVersion: 3,
		}
		if !reflect.DeepEqual(*newSessionRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GsPlayerScoresRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "8eef8847d10041e2b519ac28165e9a24.json")
			err := writeRequest(filename, []byte(`{
				"action": "groovestats/player-scores",
				"player2": {
					"chartHash": "H",
					"apiKey": "K"
				}
			}`))
			if err != nil {
				t.Fatal(err)
			}
		}()

		request := getRequest(t)
		playerScoresRequest, ok := request.(*GsPlayerScoresRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsPlayerScoresRequest{
			Id:      "8eef8847d10041e2b519ac28165e9a24",
			Player1: nil,
			Player2: &gsPlayerData{
				ChartHash: "H",
				ApiKey:    "K",
			},
		}
		if !reflect.DeepEqual(*playerScoresRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GsPlayerScoresRequestDuplicate", func(t *testing.T) {
		filename := filepath.Join(dir, "requests", "26ff392e4b894617b59c39beef4adfae.json")
		err := writeRequest(filename, []byte(`{
			"action": "groovestats/player-scores",
			"player2": {
				"chartHash": "H",
				"apiKey": "K"
			}
		}`))
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(100 * time.Millisecond)

		filename = filepath.Join(dir, "requests", "3fd78791f6a04689ad7dcb6887b418cc.json")
		err = writeRequest(filename, []byte(`{
			"action": "groovestats/player-scores",
			"player2": {
				"chartHash": "H",
				"apiKey": "K"
			}
		}`))
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(100 * time.Millisecond)

		request := getRequest(t)
		playerScoresRequest, ok := request.(*GsPlayerScoresRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsPlayerScoresRequest{
			Id:      "3fd78791f6a04689ad7dcb6887b418cc",
			Player1: nil,
			Player2: &gsPlayerData{
				ChartHash: "H",
				ApiKey:    "K",
			},
		}
		if !reflect.DeepEqual(*playerScoresRequest, expected) {
			t.Fatal("unexpected request")
		}

		if len(ipc.GsPlayerScoresRequests) > 0 {
			t.Fatal("old request not discarded")
		}
	})

	t.Run("GsPlayerLeaderboardsRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "68b50a36752141d58700b4d252dfee15.json")
			err := writeRequest(filename, []byte(`{
				"action": "groovestats/player-leaderboards",
				"player1": {
					"chartHash": "H",
					"apiKey": "K"
				}
			}`))
			if err != nil {
				t.Fatal(err)
			}
		}()

		request := getRequest(t)
		playerLeaderboardsRequest, ok := request.(*GsPlayerLeaderboardsRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsPlayerLeaderboardsRequest{
			Id: "68b50a36752141d58700b4d252dfee15",
			Player1: &gsPlayerData{
				ChartHash: "H",
				ApiKey:    "K",
			},
			Player2: nil,
		}
		if !reflect.DeepEqual(*playerLeaderboardsRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GsScoreSubmitRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "25a1506cdeff4d01b50f8207313f5db1.json")
			err := writeRequest(filename, []byte(`{
				"action": "groovestats/score-submit",
				"player1": {
					"apiKey": "K",
					"profileName": "N",
					"chartHash": "H",
					"score": 10000,
					"comment": "C675",
					"rate": 100
				},
				"player2": {
					"apiKey": "K",
					"profileName": "N",
					"chartHash": "H",
					"score": 9900,
					"comment": "C700",
					"rate": 150
				}
			}`))
			if err != nil {
				t.Fatal(err)
			}
		}()

		request := getRequest(t)
		scoreSubmitRequest, ok := request.(*GsScoreSubmitRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsScoreSubmitRequest{
			Id: "25a1506cdeff4d01b50f8207313f5db1",
			Player1: &gsScoreSubmitPlayerData{
				ApiKey:      "K",
				ProfileName: "N",
				ChartHash:   "H",
				Score:       10000,
				Comment:     "C675",
				Rate:        100,
			},
			Player2: &gsScoreSubmitPlayerData{
				ApiKey:      "K",
				ProfileName: "N",
				ChartHash:   "H",
				Score:       9900,
				Comment:     "C700",
				Rate:        150,
			},
		}
		if !reflect.DeepEqual(*scoreSubmitRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("WriteResponse", func(t *testing.T) {
		var data struct {
			Payload string `json:"payload"`
		}
		data.Payload = "foobar"

		err := ipc.WriteResponse("848431421ad846ebbb26269c749a5a43", data)
		if err != nil {
			t.Fatal(err)
		}

		filename := filepath.Join(dir, "responses", "848431421ad846ebbb26269c749a5a43.json")
		serialized, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		expected := []byte(`{"payload":"foobar"}`)
		if !bytes.Equal(serialized, expected) {
			t.Fatal("unexpected data")
		}
	})

	t.Run("Close", func(t *testing.T) {
		err := ipc.Close()
		if err != nil {
			t.Fatal(err)
		}

		_, ok := <-ipc.Requests
		if ok {
			t.Fatal("requests channel not closed")
		}
	})
}
