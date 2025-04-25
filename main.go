package main

import (
	"log"

	"github.com/OndrejHana/two-towers/server"
)

func main() {
	if err := server.Init(); err != nil {
		log.Fatal(err)
	}
}
