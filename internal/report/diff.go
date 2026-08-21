package report

import "stickerchallenge/internal/domain"

type RecordChange struct {
	ID     string
	Before domain.StickerRecord
	After  domain.StickerRecord
}

func DiffBatches(before, after domain.Batch) []RecordChange {
	changes := make([]RecordChange, 0)
	old := map[string]domain.StickerRecord{}
	for _, record := range before.Records {
		old[record.ID] = record
	}
	for _, record := range after.Records {
		previous, ok := old[record.ID]
		if !ok || !domain.CompareRecords(previous, record) {
			changes = append(changes, RecordChange{ID: record.ID, Before: previous, After: record})
		}
	}
	return changes
}
func ChangedIDs(before, after domain.Batch) []string {
	changes := DiffBatches(before, after)
	result := make([]string, 0, len(changes))
	for _, change := range changes {
		result = append(result, change.ID)
	}
	return result
}
func SameLifecycle(before, after domain.Batch) bool { return before.Status == after.Status }
func SameRecords(before, after domain.Batch) bool   { return len(DiffBatches(before, after)) == 0 }
func HasRecord(batch domain.Batch, id string) bool {
	for _, record := range batch.Records {
		if record.ID == id {
			return true
		}
	}
	return false
}
