package service

import "testing"

func TestWorkflowImportReport(t *testing.T) {
	svc := newTestService(t)
	batch, parsed, err := svc.ImportBatch("import", "Import flow", "operator", "a,21\nb,22\nc,no")
	if err != nil {
		t.Fatal(err)
	}
	if batch.ID != "import" || len(parsed.Candidates) != 2 || len(parsed.Warnings) != 1 {
		t.Fatalf("batch=%+v parsed=%+v", batch, parsed)
	}
	summary, err := svc.Summary(batch.ID)
	if err != nil || summary.Total != 2 {
		t.Fatalf("%+v %v", summary, err)
	}
}
