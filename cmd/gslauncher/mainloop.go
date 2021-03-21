package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/groovestats"
	"github.com/archiveflax/gslauncher/internal/settings"
)

func mainLoop() {
	gsClient := groovestats.NewClient()
	reload := make(chan bool)

	settings.SetUpdateCallback(func(old, new_ settings.Settings) {
		if old.SmDataDir != new_.SmDataDir {
			reload <- true
		}
	})

	for {
		smDir := settings.Get().SmDataDir
		dataDir := filepath.Join(smDir, "Save", "GrooveStats")

		info, err := os.Stat(smDir)
		if err != nil || !info.IsDir() {
			log.Print("StepMania directory not found")
			<-reload
			continue
		}

		err = os.MkdirAll(dataDir, os.ModeDir|0700)
		if err != nil {
			log.Print("failed to create data directory: ", err)
			<-reload
			continue
		}

		ipc, err := fsipc.New(dataDir)
		if err != nil {
			log.Print("failed to initialize fsipc: ", err)
			<-reload
			continue
		}

	loop:
		for {
			select {
			case request := <-ipc.Requests:
				log.Printf("REQ: %#v", request)
				processRequest(ipc, gsClient, request)
			case <-reload:
				break loop
			}
		}

		ipc.Close()
	}
}

func processRequest(ipc *fsipc.FsIpc, gsClient *groovestats.Client, request interface{}) {
	switch req := request.(type) {
	case *fsipc.PingRequest:
		response := fsipc.PingResponse{
			Version: fsipc.PingVersion{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.GsNewSessionRequest:
		resp, err := gsClient.NewSession()

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.GsPlayerScoresRequest:
		resp, err := gsClient.PlayerScores(req)

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.GsPlayerLeaderboardsRequest:
		resp, err := gsClient.PlayerLeaderboards(req)

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		ipc.WriteResponse(req.Id, response)
	case *fsipc.GsScoreSubmitRequest:
		resp, err := gsClient.ScoreSubmit(req)

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
