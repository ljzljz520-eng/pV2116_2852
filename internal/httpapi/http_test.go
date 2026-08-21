package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"stickerchallenge/internal/service"
	"stickerchallenge/internal/store"
	"testing"
)

func TestRegistrationHTTP(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	server := New(service.New(db, service.FixedClock{Value: "2116-05-01T00:00:00Z"}))
	body := []byte(`{"id":"api","label":"HTTP","owner":"op","candidates":[{"ID":"r","Number":22}]}`)
	req := httptest.NewRequest(http.MethodPost, "/batches", bytes.NewReader(body))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, req)
	if response.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/batches?label=http", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d", response.Code)
	}
}
