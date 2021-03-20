package fsipc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFsipc(t *testing.T) {
	dir := t.TempDir()

	ipc, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("PingRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "bda6a8a9d7924c149697e13b93aa68bf.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "ping"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		pingRequest, ok := request.(*PingRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := PingRequest{
			Id: "bda6a8a9d7924c149697e13b93aa68bf",
		}
		if *pingRequest != expected {
			t.Fatal("unexpected request")
		}
	})

	t.Run("NewSession", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "b787ae38ea0c465e8d853015db940915.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "groovestats/new-session"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		newSessionRequest, ok := request.(*GsNewSessionRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsNewSessionRequest{
			Id: "b787ae38ea0c465e8d853015db940915",
		}
		if *newSessionRequest != expected {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GetScoresRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "71523759f20147f79ab2b9f883492e7b.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "groovestats/get-scores",
				"api-key": "K",
				"hash": "H"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		getScoresRequest, ok := request.(*GetScoresRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GetScoresRequest{
			Id:     "71523759f20147f79ab2b9f883492e7b",
			ApiKey: "K",
			Hash:   "H",
		}
		if *getScoresRequest != expected {
			t.Fatal("unexpected request")
		}
	})

	t.Run("SubmitScoreRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "25a1506cdeff4d01b50f8207313f5db1.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "groovestats/submit-score",
				"api-key": "K",
				"profile-name": "N",
				"hash": "H",
				"score": 10000,
				"rate": 100
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		submitScoreRequest, ok := request.(*SubmitScoreRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := SubmitScoreRequest{
			Id:          "25a1506cdeff4d01b50f8207313f5db1",
			ApiKey:      "K",
			ProfileName: "N",
			Hash:        "H",
			Score:       10000,
			Rate:        100,
		}
		if *submitScoreRequest != expected {
			t.Fatal("unexpected request")
		}
	})

	t.Run("PartialWrite", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "cb95f27932174600bafab86e2e5204c7.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "ping"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}

			f, err := os.Create(filename)
			if err != nil {
				t.Fatal(err)
			}

			_, err = f.Write([]byte(`{"action": `))
			if err != nil {
				t.Fatal(err)
			}

			err = f.Sync()
			if err != nil {
				t.Fatal(err)
			}

			// give fsipc some time to process the write event
			// before issueing the next one
			time.Sleep(time.Second)

			_, err = f.Write([]byte(`"ping"}`))
			if err != nil {
				t.Fatal(err)
			}

			err = f.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		pingRequest, ok := request.(*PingRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := PingRequest{
			Id: "cb95f27932174600bafab86e2e5204c7",
		}
		if *pingRequest != expected {
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
