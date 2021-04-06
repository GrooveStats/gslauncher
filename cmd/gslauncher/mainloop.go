package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/GrooveStats/gslauncher/internal/fsipc"
	"github.com/GrooveStats/gslauncher/internal/groovestats"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

func mainLoop(unlockManager *unlocks.Manager) {
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

		var wg sync.WaitGroup

		loop := func(ch chan interface{}) {
			wg.Add(1)
			for request := range ch {
				processRequest(ipc, gsClient, request, unlockManager)
			}
			wg.Done()
		}

		go loop(ipc.GsPlayerScoresRequests)
		go loop(ipc.GsPlayerLeaderboardsRequests)
		go loop(ipc.Requests)

		<-reload
		ipc.Close()
		wg.Wait()
	}
}

func processRequest(ipc *fsipc.FsIpc, gsClient *groovestats.Client, request interface{}, unlockManager *unlocks.Manager) {
	log.Printf("REQ: %#v", request)

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
		resp, err := gsClient.NewSession(req)

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

		if err == nil {
			if req.Player1 != nil && resp.Player1 != nil && resp.Player1.Rpg != nil && resp.Player1.Rpg.Progress != nil {
				for _, quest := range resp.Player1.Rpg.Progress.QuestsCompleted {
					for _, reward := range quest.Rewards {
						if reward.SongDownloadUrl == nil {
							continue
						}
						unlockManager.AddUnlock(
							reward.Description,
							*reward.SongDownloadUrl,
							resp.Player1.Rpg.Name,
							req.Player1.ProfileName,
						)
					}
				}
			}

			if req.Player2 != nil && resp.Player2 != nil && resp.Player2.Rpg != nil && resp.Player2.Rpg.Progress != nil {
				for _, quest := range resp.Player2.Rpg.Progress.QuestsCompleted {
					for _, reward := range quest.Rewards {
						if reward.SongDownloadUrl == nil {
							continue
						}
						unlockManager.AddUnlock(
							reward.Description,
							*reward.SongDownloadUrl,
							resp.Player2.Rpg.Name,
							req.Player2.ProfileName,
						)
					}
				}
			}
		}

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		ipc.WriteResponse(req.Id, response)
	default:
		panic("unknown request type")
	}
}
