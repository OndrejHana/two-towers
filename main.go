package main

import (
	"log"
	"two-towers/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
