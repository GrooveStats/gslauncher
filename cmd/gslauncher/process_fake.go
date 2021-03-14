// +build fake

package main

import (
	"embed"
	"encoding/json"
	"math/rand"

	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/groovestats"
)

//go:embed fake/*.json
var fs embed.FS

func processRequest(ipc *fsipc.FsIpc, request interface{}) {
	switch req := request.(type) {
	case *fsipc.PingRequest:
		response := fsipc.PingResponse{Payload: req.Payload}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.GetScoresRequest:
		var filename string

		switch rand.Intn(3) {
		case 0: // network error
			response := fsipc.NetworkResponse{
				Success: false,
				Data:    (*groovestats.GetScoresResponse)(nil),
			}
			ipc.WriteResponse(req.Id, response)
			return
		case 1:
			filename = "fake/get-scores.json"
		case 2:
			filename = "fake/get-scores-rpg.json"
		}

		data, err := fs.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		var resp groovestats.GetScoresResponse
		err = json.Unmarshal(data, &resp)
		if err != nil {
			panic(err)
		}

		response := fsipc.NetworkResponse{
			Success: true,
			Data:    &resp,
		}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.SubmitScoreRequest:
		var filename string

		switch {
		case req.Rate <= 33: // network error
			response := fsipc.NetworkResponse{
				Success: false,
				Data:    (*groovestats.AutoSubmitScoreResponse)(nil),
			}
			ipc.WriteResponse(req.Id, response)
			return
		case req.Rate <= 66:
			filename = "fake/auto-submit-score-added.json"
		case req.Rate <= 100:
			filename = "fake/auto-submit-score-added-rpg.json"
		case req.Rate <= 133:
			filename = "fake/auto-submit-score-improved.json"
		case req.Rate <= 166:
			filename = "fake/auto-submit-score-not-improved.json"
		default:
			filename = "fake/auto-submit-score-not-ranked.json"
		}

		data, err := fs.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		var resp groovestats.AutoSubmitScoreResponse
		err = json.Unmarshal(data, &resp)
		if err != nil {
			panic(err)
		}

		response := fsipc.NetworkResponse{
			Success: true,
			Data:    &resp,
		}
		ipc.WriteResponse(req.Id, response)
	default:
		panic("unknown request type")
	}
}
