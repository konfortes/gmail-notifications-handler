package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	registerRoutes()

	addr := determineListenAddress()

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
	log.Printf("Listening on %s...\n", addr)
}

func registerRoutes() {
	http.HandleFunc("/health", healthHandler)
}

func determineListenAddress() string {
	port, found := os.LookupEnv("PORT")

	if !found {
		port = "3322"
		fmt.Printf("PORT env var was not supplied. defaulting to %s\n", port)
	}

	return ":" + port
}
