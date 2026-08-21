package domain

import (
	"sort"
	"strings"
)

type RecordQuery struct {
	Result    string
	Confirmed *bool
	Minimum   int
	Maximum   int
	Text      string
}
type Page struct {
	Offset int
	Limit  int
}

func MatchRecord(record StickerRecord, query RecordQuery) bool {
	if query.Result != "" && record.Result != query.Result {
		return false
	}
	if query.Confirmed != nil && record.Confirmed != *query.Confirmed {
		return false
	}
	if query.Minimum > 0 && record.Number < query.Minimum {
		return false
	}
	if query.Maximum > 0 && record.Number > query.Maximum {
		return false
	}
	if query.Text != "" && !strings.Contains(strings.ToLower(record.ID), strings.ToLower(query.Text)) {
		return false
	}
	return true
}
func FilterRecords(records []StickerRecord, query RecordQuery) []StickerRecord {
	result := make([]StickerRecord, 0)
	for _, record := range records {
		if MatchRecord(record, query) {
			result = append(result, record)
		}
	}
	SortRecords(result)
	return result
}
func PaginateRecords(records []StickerRecord, page Page) []StickerRecord {
	if page.Offset < 0 {
		page.Offset = 0
	}
	if page.Limit <= 0 {
		page.Limit = len(records)
	}
	if page.Offset >= len(records) {
		return []StickerRecord{}
	}
	end := page.Offset + page.Limit
	if end > len(records) {
		end = len(records)
	}
	return records[page.Offset:end]
}
func SortBatches(batches []Batch, descending bool) {
	sort.Slice(batches, func(i, j int) bool {
		if descending {
			return batches[i].UpdatedAt > batches[j].UpdatedAt
		}
		return batches[i].UpdatedAt < batches[j].UpdatedAt
	})
}
func Statuses() []Status {
	return []Status{StatusRegistered, StatusReviewing, StatusConfirmed, StatusPublished, StatusArchived}
}
func IsTerminal(status Status) bool   { return status == StatusArchived }
func IsReviewable(status Status) bool { return status == StatusRegistered || status == StatusReviewing }
