package report

import (
	"fmt"
	"sort"
	"stickerchallenge/internal/domain"
	"strings"
)

type RenderOptions struct {
	IncludeFailures    bool
	IncludeUnconfirmed bool
	Header             bool
}

func RenderBatch(batch domain.Batch, options RenderOptions) string {
	lines := make([]string, 0)
	if options.Header {
		lines = append(lines, "id|number|result|confirmed|divisors")
	}
	records := domain.CloneRecords(batch.Records)
	domain.SortRecords(records)
	for _, record := range records {
		if !options.IncludeFailures && record.Result == "fail" {
			continue
		}
		if !options.IncludeUnconfirmed && !record.Confirmed {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s|%d|%s|%t|%s", record.ID, record.Number, record.Result, record.Confirmed, domain.FormatDivisors(record.Divisors)))
	}
	return strings.Join(lines, "\n")
}
func RenderSummary(summary domain.Summary) string {
	return fmt.Sprintf("%s %s total=%d confirmed=%d passing=%d failing=%d", summary.BatchID, summary.Lifecycle, summary.Total, summary.Confirmed, summary.Passing, summary.Failing)
}
func RenderWarnings(warnings []string) string {
	values := append([]string(nil), warnings...)
	sort.Strings(values)
	return strings.Join(values, "\n")
}
func Exportable(batch domain.Batch) bool {
	return batch.Status == domain.StatusConfirmed || batch.Status == domain.StatusPublished || batch.Status == domain.StatusArchived
}
func ConfirmedRecords(batch domain.Batch) []domain.StickerRecord {
	result := make([]domain.StickerRecord, 0)
	for _, record := range batch.Records {
		if record.Confirmed {
			result = append(result, record)
		}
	}
	domain.SortRecords(result)
	return result
}
func PassingRecords(batch domain.Batch) []domain.StickerRecord {
	result := make([]domain.StickerRecord, 0)
	for _, record := range batch.Records {
		if record.Result == "pass" {
			result = append(result, record)
		}
	}
	domain.SortRecords(result)
	return result
}
func FailureRecords(batch domain.Batch) []domain.StickerRecord {
	result := make([]domain.StickerRecord, 0)
	for _, record := range batch.Records {
		if record.Result == "fail" {
			result = append(result, record)
		}
	}
	domain.SortRecords(result)
	return result
}
func GroupByResult(records []domain.StickerRecord) map[string][]domain.StickerRecord {
	groups := map[string][]domain.StickerRecord{"pass": {}, "fail": {}}
	for _, record := range records {
		groups[record.Result] = append(groups[record.Result], record)
	}
	return groups
}
func Labels(batches []domain.Batch) []string {
	result := make([]string, 0, len(batches))
	for _, batch := range batches {
		result = append(result, batch.Label)
	}
	sort.Strings(result)
	return result
}
