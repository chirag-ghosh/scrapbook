package server

import (
	"errors"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/", handleHealthCheck)

	log.Default().Printf("Starting server on :6969")
	err := http.ListenAndServe(":6969", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Default().Printf("Server shutdown gracefully")
	} else {
		log.Fatalf("Failed to start server: %v", err)
	}
}
