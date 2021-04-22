package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/GrooveStats/gslauncher/internal/settings"
)

func main() {
	settings.Load()

	smDir := settings.Get().SmDataDir
	dataDir := filepath.Join(smDir, "Save", "GrooveStats")

	uuid := genUuid4()

	filename := uuid + ".json"
	tmpfile := filepath.Join(dataDir, "requests", "new."+filename+".new")
	requestFile := filepath.Join(dataDir, "requests", filename)
	responseFile := filepath.Join(dataDir, "responses", filename)

	f, err := os.Create(tmpfile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, os.Stdin)
	if err != nil {
		f.Close()
		os.Remove(tmpfile)
	}

	err = os.Rename(tmpfile, requestFile)
	if err != nil {
		f.Close()
		os.Remove(tmpfile)
	}

	for i := 0; i < 60; i++ {
		time.Sleep(100 * time.Millisecond)

		f, err := os.Open(responseFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			log.Fatal(err)
		}

		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	log.Fatal("timeout")
}

func genUuid4() string {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
