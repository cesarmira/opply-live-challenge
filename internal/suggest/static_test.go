package suggest

import (
	"errors"
	"testing"
)

func TestStaticSuggest_Known(t *testing.T) {
	s := NewStatic()

	alts, err := s.Suggest("butter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alts) == 0 {
		t.Fatal("expected at least one alternative for butter")
	}
}

func TestStaticSuggest_CaseInsensitive(t *testing.T) {
	s := NewStatic()

	for _, in := range []string{"BUTTER", "  Butter  ", "butter"} {
		if _, err := s.Suggest(in); err != nil {
			t.Errorf("Suggest(%q) returned error: %v", in, err)
		}
	}
}

func TestStaticSuggest_Unknown(t *testing.T) {
	s := NewStatic()

	_, err := s.Suggest("unobtanium")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
