package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func StartServer() {
	http.HandleFunc("/timeline", handleTimeline)
	http.HandleFunc("/photo/{id}/serve", handlePhotoServe)
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/", handleHealthCheck)

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedOrigins:   []string{"http://localhost:5173"},
	})

	log.Default().Printf("Starting server on :6969")
	err := http.ListenAndServe(":6969", c.Handler(http.DefaultServeMux))
	if errors.Is(err, http.ErrServerClosed) {
		log.Default().Printf("Server shutdown gracefully")
	} else {
		log.Fatalf("Failed to start server: %v", err)
	}
}
