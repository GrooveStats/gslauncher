package main

import (
	"log"

	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/groovestats"
)

// XXX: config vars
const dataDir = "data"
const groovestatsUrl = "http://localhost:12345"

func main() {
	ipc, err := fsipc.New(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer ipc.Close()

	for request := range ipc.Requests {
		log.Printf("REQ: %#v", request)

		switch req := request.(type) {
		case *fsipc.PingRequest:
			response := fsipc.PingResponse{Payload: req.Payload}
			ipc.WriteResponse(req.Id, response)
		case *fsipc.GetScoresRequest:
			client := groovestats.NewClient(groovestatsUrl, req.ApiKey)
			resp, err := client.GetScores("somehash")

			response := fsipc.NetworkResponse{
				Success: err == nil,
				Data:    resp,
			}
			ipc.WriteResponse(req.Id, response)
		case *fsipc.SubmitScoreRequest:
			client := groovestats.NewClient(groovestatsUrl, req.ApiKey)
			resp, err := client.AutoSubmitScore("somhash", req.Rate, req.Score)

			// XXX: download unlocks

			response := fsipc.NetworkResponse{
				Success: err == nil,
				Data:    resp,
			}
			ipc.WriteResponse(req.Id, response)
		default:
			log.Fatal("unknown request")
		}
	}
}
