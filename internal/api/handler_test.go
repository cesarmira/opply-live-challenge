package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cesarmirasanchez/opply-live-challenge/internal/suggest"
)

func newTestServer() http.Handler {
	return NewMux(suggest.NewStatic())
}

func TestSuggest_OK(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/suggest?ingredient=butter", nil)

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body response
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body.Ingredient != "butter" || len(body.Alternatives) == 0 {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestSuggest_NotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/suggest?ingredient=unobtanium", nil)

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestSuggest_MissingIngredient(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/suggest", nil)

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
