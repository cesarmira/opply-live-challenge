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
	apiKey := os.Getenv("OPENCODE_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENCODE_API_KEY is required")
	}
	baseURL := os.Getenv("OPENCODE_BASE_URL")
	if baseURL == "" {
		log.Fatal("OPENCODE_BASE_URL is required")
	}

	addr := ":" + port()
	mux := api.NewMux(suggest.NewLLM(baseURL, apiKey, "deepseek-v4-flash"))

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
