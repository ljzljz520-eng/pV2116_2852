package importer

import "testing"

func TestParseRows(t *testing.T) {
	result, err := ParseRows("a,21\nb,no\nmissing\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Candidates) != 1 || result.Candidates[0].Number != 21 {
		t.Fatalf("%+v", result)
	}
	if len(result.Warnings) != 2 {
		t.Fatalf("warnings=%v", result.Warnings)
	}
}

func TestParseJSON(t *testing.T) {
	result, err := ParseJSON([]byte(`[{"ID":"b","Number":7},{"ID":"a","Number":5}]`))
	if err != nil || result.Candidates[0].ID != "a" {
		t.Fatalf("%+v %v", result, err)
	}
}
