package main

import (
	"log"
	"two-towers/backend"
)

func main() {
	if err := backend.Serve(); err != nil {
		log.Fatal("could not run the server")
		log.Fatal(err)
	}
}
