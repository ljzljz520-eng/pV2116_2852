package store

import (
	"path/filepath"
	"stickerchallenge/internal/domain"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "challenge.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	batch := domain.Batch{ID: "b", Label: "label", Owner: "owner", Status: domain.StatusRegistered, Version: 1, Records: []domain.StickerRecord{{ID: "r", BatchID: "b", Number: 2, Result: "pass"}}}
	if err := s.PutBatch(batch); err != nil {
		t.Fatal(err)
	}
	if err := s.PutEvent(domain.AuditEvent{ID: "e", BatchID: "b", Action: "register"}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutNote(domain.CollaborationNote{ID: "n", BatchID: "b", Author: "a", Body: "body"}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutExport(domain.ExportSnapshot{ID: "x", BatchID: "b", Payload: "[]"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if got, err := s.GetBatch("b"); err != nil || got.Label != "label" {
		t.Fatalf("batch=%+v err=%v", got, err)
	}
	if got, err := s.ListEvents("b"); err != nil || len(got) != 1 {
		t.Fatalf("events=%v err=%v", got, err)
	}
	if got, err := s.ListNotes("b"); err != nil || len(got) != 1 {
		t.Fatalf("notes=%v err=%v", got, err)
	}
	if got, err := s.GetExport("x"); err != nil || got.Payload != "[]" {
		t.Fatalf("export=%+v err=%v", got, err)
	}
}

func TestVersionConflict(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "challenge.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	batch := domain.Batch{ID: "b", Label: "label", Owner: "owner", Status: domain.StatusRegistered, Version: 1, Records: []domain.StickerRecord{{ID: "r", BatchID: "b", Number: 2}}}
	if err := s.PutBatch(batch); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateBatch(batch, 2); err != domain.ErrConflict {
		t.Fatalf("err=%v", err)
	}
}
