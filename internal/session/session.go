package session

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/GrooveStats/gslauncher/internal/fsipc"
	"github.com/GrooveStats/gslauncher/internal/groovestats"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

type Session struct {
	unlockManager *unlocks.Manager
	gsClient      *groovestats.Client
	ipc           *fsipc.FsIpc
	cmd           *exec.Cmd
	wg            sync.WaitGroup
}

func Launch(unlockManager *unlocks.Manager) (*Session, error) {
	sess := &Session{
		unlockManager: unlockManager,
		gsClient:      groovestats.NewClient(),
	}

	err := sess.startIpc()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize fsipc: %w", err)
	}

	err = sess.startSM()
	if err != nil {
		sess.ipc.Close()
		sess.wg.Wait()
		return nil, fmt.Errorf("failed to run StepMania: %w", err)
	}

	go func() {
		sess.wg.Add(1)
		sess.cmd.Wait()
		sess.ipc.Close()
		sess.wg.Done()
	}()

	return sess, nil
}

func (sess *Session) Wait() {
	sess.wg.Wait()
}

func (sess *Session) Kill() {
	sess.cmd.Process.Kill()
	sess.wg.Wait()
}

func (sess *Session) startIpc() error {
	smDir := settings.Get().SmDataDir
	dataDir := filepath.Join(smDir, "Save", "GrooveStats")

	_, err := os.Stat(smDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dataDir, os.ModeDir|0700)
	if err != nil {
		return err
	}

	ipc, err := fsipc.New(dataDir)
	if err != nil {
		return err
	}

	processChannel := func(ch chan interface{}) {
		sess.wg.Add(1)
		for request := range ch {
			sess.processRequest(request)
		}
		sess.wg.Done()
	}

	go processChannel(ipc.GsPlayerScoresRequests)
	go processChannel(ipc.GsPlayerLeaderboardsRequests)
	go processChannel(ipc.Requests)

	sess.ipc = ipc
	return nil
}

func (sess *Session) startSM() error {
	cmd := exec.Command(settings.Get().SmExePath)

	err := cmd.Start()
	if err != nil {
		return err
	}

	sess.cmd = cmd
	return nil
}

func (sess *Session) processRequest(request interface{}) {
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
		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsNewSessionRequest:
		resp, err := sess.gsClient.NewSession(req)

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsPlayerScoresRequest:
		resp, err := sess.gsClient.PlayerScores(req)

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsPlayerLeaderboardsRequest:
		resp, err := sess.gsClient.PlayerLeaderboards(req)

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsScoreSubmitRequest:
		resp, err := sess.gsClient.ScoreSubmit(req)

		if err == nil {
			if req.Player1 != nil && resp.Player1 != nil && resp.Player1.Rpg != nil && resp.Player1.Rpg.Progress != nil {
				for _, quest := range resp.Player1.Rpg.Progress.QuestsCompleted {
					if quest.SongDownloadUrl == nil {
						continue
					}

					descriptions := make([]string, 0)
					for _, reward := range quest.Rewards {
						if reward.Type == "song" {
							descriptions = append(descriptions, reward.Description)
						}
					}

					sess.unlockManager.AddUnlock(
						quest.Title,
						*quest.SongDownloadUrl,
						resp.Player1.Rpg.Name,
						req.Player1.ProfileName,
						descriptions,
					)
				}
			}

			if req.Player2 != nil && resp.Player2 != nil && resp.Player2.Rpg != nil && resp.Player2.Rpg.Progress != nil {
				for _, quest := range resp.Player2.Rpg.Progress.QuestsCompleted {
					if quest.SongDownloadUrl == nil {
						continue
					}

					descriptions := make([]string, 0)
					for _, reward := range quest.Rewards {
						if reward.Type == "song" {
							descriptions = append(descriptions, reward.Description)
						}
					}

					sess.unlockManager.AddUnlock(
						quest.Title,
						*quest.SongDownloadUrl,
						resp.Player2.Rpg.Name,
						req.Player2.ProfileName,
						descriptions,
					)
				}
			}
		}

		response := fsipc.NetworkResponse{
			Success: err == nil,
			Data:    resp,
		}
		sess.ipc.WriteResponse(req.Id, response)
	default:
		panic("unknown request type")
	}
}
