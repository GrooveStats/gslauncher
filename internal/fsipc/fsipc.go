package fsipc

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
)

type FsIpc struct {
	Requests                     chan interface{}
	GsPlayerScoresRequests       chan interface{}
	GsPlayerLeaderboardsRequests chan interface{}
	RootDir                      string

	requestDir  string
	responseDir string
	watcher     *fsnotify.Watcher
	shutdown    chan bool
	wg          sync.WaitGroup
}

func New(rootDir string) (*FsIpc, error) {
	info, err := os.Stat(rootDir)
	if os.IsNotExist(err) || !info.IsDir() {
		return nil, fmt.Errorf("root directory doesn't exist")
	}

	requestDir := filepath.Join(rootDir, "requests")
	err = os.RemoveAll(requestDir)
	if err != nil {
		return nil, err
	}
	err = os.Mkdir(requestDir, os.ModeDir|0700)
	if err != nil {
		return nil, err
	}

	responseDir := filepath.Join(rootDir, "responses")
	err = os.RemoveAll(responseDir)
	if err != nil {
		return nil, err
	}
	err = os.Mkdir(responseDir, os.ModeDir|0700)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fsipc := FsIpc{
		Requests:                     make(chan interface{}, 5),
		GsPlayerScoresRequests:       make(chan interface{}, 1),
		GsPlayerLeaderboardsRequests: make(chan interface{}, 1),
		RootDir:                      rootDir,
		requestDir:                   requestDir,
		responseDir:                  responseDir,
		watcher:                      watcher,
		shutdown:                     make(chan bool),
	}

	err = fsipc.watcher.Add(fsipc.requestDir)
	if err != nil {
		fsipc.watcher.Close()
		return nil, err
	}

	fsipc.wg.Add(1)
	go fsipc.loop()

	return &fsipc, nil
}

func (fsipc *FsIpc) Close() error {
	close(fsipc.shutdown)
	fsipc.wg.Wait()

	close(fsipc.Requests)
	close(fsipc.GsPlayerScoresRequests)
	close(fsipc.GsPlayerLeaderboardsRequests)

	err := os.RemoveAll(fsipc.requestDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(fsipc.responseDir)
	if err != nil {
		return err
	}

	return fsipc.watcher.Close()
}

func (fsipc *FsIpc) loop() {
	defer fsipc.wg.Done()

	for {
		select {
		case event, ok := <-fsipc.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				fsipc.wg.Add(1)
				go fsipc.handleFile(event.Name)
			}
		case err, ok := <-fsipc.watcher.Errors:
			if !ok {
				return
			}

			log.Print("fsnotify error: ", err)
		case <-fsipc.shutdown:
			return
		}
	}
}

func readFilePatient(filename string) ([]byte, error) {
	var data []byte
	var err error

	if runtime.GOOS == "windows" {
		for retry := 0; retry < 30; retry++ {
			data, err = os.ReadFile(filename)

			if patherr, ok := err.(*fs.PathError); ok {
				// ERROR_SHARING_VIOLATION
				if patherr.Err == syscall.Errno(0x20) {
					time.Sleep(100 * time.Millisecond)
					continue
				}
			}

			break
		}
	} else {
		data, err = os.ReadFile(filename)
	}

	return data, err
}

func (fsipc *FsIpc) handleFile(filename string) {
	defer fsipc.wg.Done()

	if !strings.HasSuffix(filename, ".json") {
		return
	}

	basename := filepath.Base(filename)
	id := strings.TrimSuffix(basename, ".json")

	info, err := os.Stat(filename)
	if err != nil {
		log.Printf("failed to stat %s: %v", basename, err)
		return
	}

	if !info.Mode().IsRegular() {
		return
	}

	// SM only waits up to one minute for a reply, so if the request is too
	// old, just discard it.
	if info.ModTime().Add(time.Minute).Before(time.Now()) {
		log.Print("discarding stale request: ", id)
		err = os.Remove(filename)
		if err != nil {
			log.Printf("failed to delete %s: %v", basename, err)
		}
		return
	}

	data, err := readFilePatient(filename)
	if err != nil {
		log.Printf("failed to read %s: %v", basename, err)
		return
	}

	var base struct {
		Action string `json:"action"`
	}

	err = json.Unmarshal(data, &base)
	if err != nil {
		log.Printf("failed to unmarshal request %s: %v", id, err)
		return
	}

	var request interface{}

	switch base.Action {
	case "ping":
		request = &PingRequest{Id: id}
	case "groovestats/new-session":
		request = &GsNewSessionRequest{Id: id}
	case "groovestats/player-scores":
		request = &GsPlayerScoresRequest{Id: id}
	case "groovestats/player-leaderboards":
		request = &GsPlayerLeaderboardsRequest{Id: id}
	case "groovestats/score-submit":
		request = &GsScoreSubmitRequest{Id: id}
	case "":
		log.Printf("invalid request %s: missing action", id)
		return
	default:
		log.Printf("invalid request %s: unknown action %s", id, base.Action)
		return
	}

	err = json.Unmarshal(data, request)
	if err != nil {
		log.Printf("failed to unmarshal request %s: %v", id, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Printf("invalid request %s: %v", id, err)
		return
	}

	switch req := request.(type) {
	case *GsPlayerScoresRequest:
		// empty the buffer
		select {
		case <-fsipc.GsPlayerScoresRequests:
			// do nothing
		default:
			// do nothing
		}
		fsipc.GsPlayerScoresRequests <- req
	case *GsPlayerLeaderboardsRequest:
		// empty the buffer
		select {
		case <-fsipc.GsPlayerLeaderboardsRequests:
			// do nothing
		default:
			// do nothing
		}
		fsipc.GsPlayerLeaderboardsRequests <- req
	default:
		fsipc.Requests <- request
	}

	err = os.Remove(filename)
	if err != nil {
		log.Printf("failed to delete %s: %v", basename, err)
	}

	return
}

func (fsipc *FsIpc) WriteResponse(id string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	filename := filepath.Join(fsipc.responseDir, id+".json")
	err = os.WriteFile(filename, b, 0600)

	if err == nil {
		// SM only waits up to one minute for a reply, so when the
		// response is too old, discard it.
		fsipc.wg.Add(1)
		go func() {
			defer fsipc.wg.Done()

			select {
			case <-time.After(time.Minute):
				err := os.Remove(filename)
				if err != nil {
					log.Printf("failed to delete %s: %v", filename, err)
				}
			case <-fsipc.shutdown:
				return
			}
		}()
	}

	return err
}
