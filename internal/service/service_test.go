package service

import (
	"path/filepath"
	"stickerchallenge/internal/domain"
	"stickerchallenge/internal/store"
	"testing"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return New(db, FixedClock{Value: "2116-05-01T00:00:00Z"})
}

func TestRegisterAndSearch(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.RegisterBatch("b", "Boundary", "operator", []domain.Candidate{{ID: "r", Number: 22}}); err != nil {
		t.Fatal(err)
	}
	items, err := svc.Search(domain.SearchQuery{Label: "bound"})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%v err=%v", items, err)
	}
}
