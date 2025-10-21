package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"moonbeam/internal"
)

func TestHelloEndpoint(t *testing.T) {
	e := internal.NewRouter()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if got := rec.Body.String(); len(got) == 0 {
		t.Fatalf("expected non‑empty response body")
	}
}
