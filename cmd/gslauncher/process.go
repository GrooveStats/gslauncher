// +build !fake

package main

import (
	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/groovestats"
)

func processRequest(ipc *fsipc.FsIpc, request interface{}) {
	switch req := request.(type) {
	case *fsipc.PingRequest:
		response := fsipc.PingResponse{Payload: req.Payload}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.GetScoresRequest:
		client := groovestats.NewClient(groovestatsUrl, req.ApiKey)
		resp, err := client.GetScores(req.Hash)

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.SubmitScoreRequest:
		client := groovestats.NewClient(groovestatsUrl, req.ApiKey)
		resp, err := client.AutoSubmitScore(req.Hash, req.Rate, req.Score)

		// XXX: download unlocks

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		ipc.WriteResponse(req.Id, response)
	default:
		panic("unknown request type")
	}
}
