package service

import (
	"stickerchallenge/internal/domain"
	"testing"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	svc := newTestService(t)
	batch, err := svc.RegisterBatch("search", "Search flow", "operator", []domain.Candidate{{ID: "r1", Number: 22}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.StartReview(batch.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ConfirmBatch(batch.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.UpdateRecord(batch.ID, "r1", "operator", 25, 3); err != nil {
		t.Fatal(err)
	}
	updated, err := svc.Store.GetBatch(batch.ID)
	if err != nil || updated.Records[0].Result != "pass" {
		t.Fatalf("%+v %v", updated, err)
	}
	if _, err = svc.PublishBatch(batch.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	snapshot, err := svc.ExportConfirmed(batch.ID, "operator")
	if err != nil || snapshot.Payload == "" {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
}
