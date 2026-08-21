package service

import (
	"stickerchallenge/internal/domain"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	svc := newTestService(t)
	batch, err := svc.RegisterBatch("create", "Create flow", "operator", []domain.Candidate{{ID: "r1", Number: 30}, {ID: "r2", Number: 70}})
	if err != nil {
		t.Fatal(err)
	}
	if batch.Status != domain.StatusRegistered {
		t.Fatal(batch.Status)
	}
	if _, err = svc.StartReview(batch.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ConfirmBatch(batch.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.PublishBatch(batch.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	archived, err := svc.ArchiveBatch(batch.ID, "auditor")
	if err != nil || archived.Status != domain.StatusArchived {
		t.Fatalf("%+v %v", archived, err)
	}
}
