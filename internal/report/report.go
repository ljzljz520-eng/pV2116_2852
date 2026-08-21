package report

import (
	"encoding/csv"
	"encoding/json"
	"sort"
	"stickerchallenge/internal/domain"
	"strconv"
	"strings"
)

func FilterBatches(batches []domain.Batch, query domain.SearchQuery) []domain.Batch {
	result := make([]domain.Batch, 0)
	for _, batch := range batches {
		if query.Label != "" && !strings.Contains(strings.ToLower(batch.Label), strings.ToLower(query.Label)) {
			continue
		}
		if query.Status != "" && batch.Status != query.Status {
			continue
		}
		result = append(result, domain.CloneBatch(batch))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func Summarize(batch domain.Batch) domain.Summary {
	passing, failing, confirmed := domain.CountResults(batch.Records)
	return domain.Summary{BatchID: batch.ID, Label: batch.Label, Total: len(batch.Records), Confirmed: confirmed, Passing: passing, Failing: failing, Lifecycle: batch.Status}
}

func JSONSummary(batch domain.Batch) ([]byte, error) { return json.Marshal(Summarize(batch)) }

func CSVSummary(batch domain.Batch) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"id", "number", "result", "confirmed"}); err != nil {
		return "", err
	}
	for _, record := range batch.Records {
		if err := writer.Write([]string{record.ID, strconv.Itoa(record.Number), record.Result, strconv.FormatBool(record.Confirmed)}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return builder.String(), writer.Error()
}
