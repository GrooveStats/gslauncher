package main

import (
	"log"

	"github.com/archiveflax/gslauncher/internal/fsipc"
)

func main() {
	ipc, err := fsipc.New("data")
	if err != nil {
		log.Fatal(err)
	}
	defer ipc.Close()

	for request := range ipc.Requests {
		log.Printf("REQ: %#v", request)
	}
}
