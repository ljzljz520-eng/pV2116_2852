package service

import (
	"stickerchallenge/internal/domain"
	"strings"
	"testing"
)

func TestBusiness05Regression(t *testing.T) {
	svc := newTestService(t)
	batch, err := svc.RegisterBatch("2116-05", "Boundary batch", "operator", []domain.Candidate{{ID: "r1", Number: 22}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.StartReview(batch.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ConfirmBatch(batch.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ExportConfirmed(batch.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.UpdateRecord(batch.ID, "r1", "operator", 25, 3); err != nil {
		t.Fatal(err)
	}
	second, err := svc.ExportConfirmed(batch.ID, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(second.Payload, `"number":22`) {
		t.Fatalf("stale export retained old number: %s", second.Payload)
	}
}
