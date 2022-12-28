package main

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	//
	// Initialize the router
	//
	router := chi.NewRouter()
	//
	// Start the server
	//
	app := &http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: router,
	}
	log.Fatal().Err(app.ListenAndServe()).Msg("starting the http service")
}
