package report

import "stickerchallenge/internal/domain"

type Metrics struct {
	Total      int
	Passing    int
	Failing    int
	Confirmed  int
	Completion float64
}

func CalculateMetrics(batch domain.Batch) Metrics {
	summary := Summarize(batch)
	completion := float64(0)
	if summary.Total > 0 {
		completion = float64(summary.Confirmed) / float64(summary.Total)
	}
	return Metrics{Total: summary.Total, Passing: summary.Passing, Failing: summary.Failing, Confirmed: summary.Confirmed, Completion: completion}
}
func MergeMetrics(values []Metrics) Metrics {
	result := Metrics{}
	for _, value := range values {
		result.Total += value.Total
		result.Passing += value.Passing
		result.Failing += value.Failing
		result.Confirmed += value.Confirmed
	}
	if result.Total > 0 {
		result.Completion = float64(result.Confirmed) / float64(result.Total)
	}
	return result
}
func StatusCounts(batches []domain.Batch) map[domain.Status]int {
	counts := map[domain.Status]int{}
	for _, batch := range batches {
		counts[batch.Status]++
	}
	return counts
}
func ResultCounts(records []domain.StickerRecord) map[string]int {
	counts := map[string]int{}
	for _, record := range records {
		counts[record.Result]++
	}
	return counts
}
func CompletionLabel(metrics Metrics) string {
	if metrics.Completion >= 1 {
		return "complete"
	}
	if metrics.Completion > 0 {
		return "partial"
	}
	return "pending"
}
func HasFailures(metrics Metrics) bool { return metrics.Failing > 0 }
