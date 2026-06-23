// Command server starts the ingredient-substitution HTTP API.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cesarmirasanchez/opply-live-challenge/internal/api"
	"github.com/cesarmirasanchez/opply-live-challenge/internal/suggest"
)

func main() {
	addr := ":" + port()

	mux := api.NewMux(suggest.NewStatic())

	log.Printf("ingredient-substitution API listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

// port returns the listen port from PORT, defaulting to 8080.
func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
