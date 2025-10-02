package main

import (
	"log"

	"bookcover-api/internal/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
