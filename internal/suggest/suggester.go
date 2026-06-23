// Package suggest defines the core domain for ingredient substitution.
//
// The Suggester interface is the pluggable contract: the HTTP layer depends
// only on this interface, so a static map, a database, or an LLM-backed
// implementation can be swapped in without touching the transport code.
package suggest

import "errors"

// ErrNotFound is returned when no alternatives are known for an ingredient.
var ErrNotFound = errors.New("no alternatives found")

// Alternative is a single suggested substitute for an ingredient.
type Alternative struct {
	Name  string `json:"name"`
	Notes string `json:"notes,omitempty"`
}

// Suggester returns alternative ingredients for a given ingredient.
//
// Implementations should return ErrNotFound when the ingredient is unknown,
// so callers can distinguish "no data" from real failures.
type Suggester interface {
	Suggest(ingredient string) ([]Alternative, error)
}
