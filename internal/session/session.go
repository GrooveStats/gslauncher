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
	"github.com/GrooveStats/gslauncher/internal/version"
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

	sess.wg.Add(1)
	go func() {
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
		defer sess.wg.Done()

		for request := range ch {
			sess.processRequest(request)
		}
	}

	sess.wg.Add(3)
	go processChannel(ipc.GsPlayerScoresRequests)
	go processChannel(ipc.GsPlayerLeaderboardsRequests)
	go processChannel(ipc.Requests)

	sess.ipc = ipc
	return nil
}

func (sess *Session) startSM() error {
	smExePath := settings.Get().SmExePath

	// Let's launch StepMania! We also have to set the working directory,
	// because SM 5.3 (outfox) for Linux searches for bundled shared
	// libraries in the current working directory.
	cmd := exec.Command(smExePath)
	cmd.Dir = filepath.Dir(smExePath)

	err := cmd.Start()
	if err != nil {
		return err
	}

	sess.cmd = cmd
	return nil
}

func newNetworkResponse(data interface{}, err error) *fsipc.NetworkResponse {
	status := "success"
	if err != nil {
		switch err.(type) {
		case *groovestats.DisabledError:
			status = "disabled"
		default:
			status = "fail"
		}
	}

	return &fsipc.NetworkResponse{
		Status: status,
		Data:   data,
	}
}

func (sess *Session) processRequest(request interface{}) {
	switch req := request.(type) {
	case *fsipc.PingRequest:
		if req.Protocol != version.Protocol {
			break
		}

		response := fsipc.PingResponse{
			Version: fsipc.PingVersion{
				Major: version.Major,
				Minor: version.Minor,
				Patch: version.Patch,
			},
		}

		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsNewSessionRequest:
		resp, err := sess.gsClient.NewSession(req)

		if err != nil {
			log.Printf("failed to start new session: %v", err)
		}

		response := newNetworkResponse(resp, err)
		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsPlayerScoresRequest:
		resp, err := sess.gsClient.PlayerScores(req)

		if err != nil {
			log.Printf("failed to fetch player scores: %v", err)
		}

		response := newNetworkResponse(resp, err)
		sess.ipc.WriteResponse(req.Id, response)
	case *fsipc.GsPlayerLeaderboardsRequest:
		resp, err := sess.gsClient.PlayerLeaderboards(req)

		if err != nil {
			log.Printf("failed to fetch player leaderboards: %v", err)
		}

		response := newNetworkResponse(resp, err)
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
		} else {
			log.Printf("failed to submit score: %v", err)
		}

		response := newNetworkResponse(resp, err)
		sess.ipc.WriteResponse(req.Id, response)
	default:
		panic("unknown request type")
	}
}
