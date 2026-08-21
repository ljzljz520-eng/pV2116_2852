package report

import (
	"stickerchallenge/internal/domain"
	"strings"
)

func ContainsLabel(batch domain.Batch, text string) bool {
	return text == "" || strings.Contains(strings.ToLower(batch.Label), strings.ToLower(text))
}
func ByOwner(batches []domain.Batch, owner string) []domain.Batch {
	result := make([]domain.Batch, 0)
	for _, batch := range batches {
		if batch.Owner == owner {
			result = append(result, batch)
		}
	}
	return result
}
func ByStatus(batches []domain.Batch, status domain.Status) []domain.Batch {
	result := make([]domain.Batch, 0)
	for _, batch := range batches {
		if batch.Status == status {
			result = append(result, batch)
		}
	}
	return result
}
func WithRecords(batches []domain.Batch) []domain.Batch {
	result := make([]domain.Batch, 0)
	for _, batch := range batches {
		if len(batch.Records) > 0 {
			result = append(result, batch)
		}
	}
	return result
}
func RecordIDs(batch domain.Batch) []string {
	result := make([]string, 0, len(batch.Records))
	for _, record := range batch.Records {
		result = append(result, record.ID)
	}
	return result
}
func NumberRange(batch domain.Batch) (int, int) {
	if len(batch.Records) == 0 {
		return 0, 0
	}
	min, max := batch.Records[0].Number, batch.Records[0].Number
	for _, record := range batch.Records[1:] {
		if record.Number < min {
			min = record.Number
		}
		if record.Number > max {
			max = record.Number
		}
	}
	return min, max
}
