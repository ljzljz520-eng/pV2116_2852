package domain

import "testing"

func TestEvaluateNumber(t *testing.T) {
	result, hits, err := EvaluateNumber(210, []int{2, 3, 5, 7})
	if err != nil || result != "pass" || len(hits) != 4 {
		t.Fatalf("result=%q hits=%v err=%v", result, hits, err)
	}
	if _, _, err := EvaluateNumber(0, []int{2}); err == nil {
		t.Fatal("expected invalid number")
	}
}

func TestLifecycle(t *testing.T) {
	batch := Batch{Status: StatusRegistered, Version: 1}
	if err := Transition(&batch, StatusReviewing); err != nil {
		t.Fatal(err)
	}
	if err := Transition(&batch, StatusArchived); err == nil {
		t.Fatal("expected rejected transition")
	}
}
