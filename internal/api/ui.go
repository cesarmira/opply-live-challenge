package api

import (
	_ "embed"
	"net/http"
)

//go:embed static/index.html
var indexHTML []byte

func uiHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}
