package fsipc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
)

type FsIpc struct {
	Requests chan interface{}
	RootDir  string

	requestDir    string
	responseDir   string
	watcher       *fsnotify.Watcher
	cleanupTicker *time.Ticker
	shutdown      chan bool
}

func New(rootDir string) (*FsIpc, error) {
	info, err := os.Stat(rootDir)
	if os.IsNotExist(err) || !info.IsDir() {
		return nil, fmt.Errorf("root directory doesn't exist")
	}

	requestDir := filepath.Join(rootDir, "requests")
	info, err = os.Stat(requestDir)
	if os.IsNotExist(err) || !info.IsDir() {
		err = os.Mkdir(requestDir, os.ModeDir|0700)
		if err != nil {
			return nil, err
		}
	}

	responseDir := filepath.Join(rootDir, "responses")
	info, err = os.Stat(responseDir)
	if os.IsNotExist(err) || !info.IsDir() {
		err = os.Mkdir(responseDir, os.ModeDir|0700)
		if err != nil {
			return nil, err
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fsipc := FsIpc{
		Requests:      make(chan interface{}),
		RootDir:       rootDir,
		requestDir:    requestDir,
		responseDir:   responseDir,
		watcher:       watcher,
		cleanupTicker: time.NewTicker(time.Minute),
		shutdown:      make(chan bool),
	}

	err = fsipc.watcher.Add(fsipc.requestDir)
	if err != nil {
		fsipc.watcher.Close()
		return nil, err
	}

	go fsipc.loop()

	return &fsipc, nil
}

func (fsipc *FsIpc) Close() error {
	fsipc.cleanupTicker.Stop()
	fsipc.shutdown <- true

	close(fsipc.Requests)
	close(fsipc.shutdown)

	return fsipc.watcher.Close()
}

func (fsipc *FsIpc) loop() {
	entries, err := os.ReadDir(fsipc.requestDir)
	if err == nil {
		for _, entry := range entries {
			fsipc.handleFile(filepath.Join(fsipc.requestDir, entry.Name()))
		}
	} else {
		log.Print("failed to list requests: ", err)
	}

	fsipc.cleanup()

loop:
	for {
		select {
		case event, ok := <-fsipc.watcher.Events:
			if !ok {
				break loop
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				fsipc.handleFile(event.Name)
			}
		case err, ok := <-fsipc.watcher.Errors:
			if !ok {
				break loop
			}

			log.Print("fsnotify error: ", err)
		case <-fsipc.cleanupTicker.C:
			fsipc.cleanup()
		case <-fsipc.shutdown:
			break loop
		}
	}
}

func (fsipc *FsIpc) handleFile(filename string) {
	basename := filepath.Base(filename)
	id := strings.TrimSuffix(basename, ".json")

	info, err := os.Stat(filename)
	if err != nil {
		log.Printf("failed to stat %s: %v", basename, err)
		return
	}

	if !info.Mode().IsRegular() || !strings.HasSuffix(filename, ".json") {
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

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read %s: %v", basename, err)
		return
	}

	if !json.Valid(data) {
		// If we get here the request file is probaly not fully written
		// yet. Let's try again on the next write event.
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
	case "groovestats/get-scores":
		request = &GetScoresRequest{Id: id}
	case "groovestats/submit-score":
		request = &SubmitScoreRequest{Id: id}
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

	fsipc.Requests <- request

	err = os.Remove(filename)
	if err != nil {
		log.Printf("failed to delete %s: %v", basename, err)
	}

	return
}

func (fsipc *FsIpc) cleanup() {
	entries, err := os.ReadDir(fsipc.responseDir)
	if err != nil {
		log.Print("failed to list responses: ", err)
		return
	}

	for _, entry := range entries {
		filename := filepath.Join(fsipc.responseDir, entry.Name())

		info, err := os.Stat(filename)
		if err != nil {
			log.Printf("failed to stat %s: %v", entry.Name(), err)
			continue
		}

		if !info.Mode().IsRegular() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// SM only waits up to one minute for a reply, so if the
		// response is too old, discard it.
		if info.ModTime().Add(time.Minute).Before(time.Now()) {
			err = os.Remove(filename)
			if err != nil {
				log.Printf("failed to delete %s: %v", entry.Name(), err)
			}
		}
	}
}

func (fsipc *FsIpc) WriteResponse(id string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	filename := filepath.Join(fsipc.responseDir, id+".json")
	return os.WriteFile(filename, b, 0600)
}
