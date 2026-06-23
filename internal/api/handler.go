// Package api is the HTTP transport layer. It translates JSON requests into
// calls on a suggest.Suggester and writes JSON responses. It contains no
// substitution logic of its own.
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cesarmirasanchez/opply-live-challenge/internal/suggest"
)

// request is the JSON body accepted by POST /suggest.
type request struct {
	Ingredient string `json:"ingredient"`
}

// response is the JSON body returned by POST /suggest.
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
	mux.HandleFunc("POST /suggest", suggestHandler(s))
	mux.HandleFunc("GET /healthz", healthHandler)
	return mux
}

// suggestHandler returns a handler that suggests alternatives for an ingredient.
func suggestHandler(s suggest.Suggester) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if req.Ingredient == "" {
			writeError(w, http.StatusBadRequest, "field \"ingredient\" is required")
			return
		}

		alts, err := s.Suggest(req.Ingredient)
		if errors.Is(err, suggest.ErrNotFound) {
			writeError(w, http.StatusNotFound, "no alternatives found for \""+req.Ingredient+"\"")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, http.StatusOK, response{Ingredient: req.Ingredient, Alternatives: alts})
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
