// Package api is the HTTP transport layer. It reads the ingredient from the
// request and writes JSON responses by calling a suggest.Suggester. It contains
// no substitution logic of its own.
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cesarmirasanchez/opply-live-challenge/internal/suggest"
)

// response is the JSON body returned by GET /suggest.
type response struct {
	Ingredient   string                `json:"ingredient"`
	Alternatives []suggest.Alternative `json:"alternatives"`
}

// errorBody is the JSON body returned for any error.
type errorBody struct {
	Error string `json:"error"`
}

// NewMux builds the HTTP routes wired to the given Suggester.
func NewMux(s suggest.Suggester) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /suggest", suggestHandler(s))
	mux.HandleFunc("GET /healthz", healthHandler)
	return mux
}

// suggestHandler returns a handler that suggests alternatives for an ingredient.
func suggestHandler(s suggest.Suggester) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredient := r.URL.Query().Get("ingredient")
		if ingredient == "" {
			writeError(w, http.StatusBadRequest, "query parameter \"ingredient\" is required")
			return
		}

		alts, err := s.Suggest(ingredient)
		if errors.Is(err, suggest.ErrNotFound) {
			writeError(w, http.StatusNotFound, "no alternatives found for \""+ingredient+"\"")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, http.StatusOK, response{Ingredient: ingredient, Alternatives: alts})
	}
}

// healthHandler is a trivial liveness check used by the smoke test.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorBody{Error: msg})
}
