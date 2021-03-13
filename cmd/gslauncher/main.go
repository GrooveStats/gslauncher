package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/groovestats"
	"github.com/archiveflax/gslauncher/internal/gui"
	"github.com/archiveflax/gslauncher/internal/settings"
)

const groovestatsUrl = "http://localhost:12345" // XXX

func main() {
	settings.Load()

	go mainLoop()

	app := gui.NewApp()
	app.Run()
}

func mainLoop() {
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
					log.Print("unknown request")
				}
			case <-reload:
				break loop
			}
		}

		ipc.Close()
	}
}
