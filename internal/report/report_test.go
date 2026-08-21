package report

import (
	"stickerchallenge/internal/domain"
	"strings"
	"testing"
)

func TestFilterAndSummary(t *testing.T) {
	batches := []domain.Batch{{ID: "b", Label: "Boundary", Status: domain.StatusConfirmed, Records: []domain.StickerRecord{{Result: "pass", Confirmed: true}, {Result: "fail"}}}}
	got := FilterBatches(batches, domain.SearchQuery{Label: "bound"})
	if len(got) != 1 {
		t.Fatal("expected batch")
	}
	summary := Summarize(got[0])
	if summary.Passing != 1 || summary.Failing != 1 || summary.Confirmed != 1 {
		t.Fatalf("%+v", summary)
	}
	csv, err := CSVSummary(got[0])
	if err != nil || !strings.Contains(csv, "confirmed") {
		t.Fatalf("%q %v", csv, err)
	}
}
