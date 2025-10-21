package main

import (
	"log"
	"moonbeam/internal"
)

func main() {
	e := internal.NewRouter()

	log.Println("Starting internal on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("internal error: %v", err)
	}
}
