package domain

import (
	"fmt"
	"sort"
	"strings"
)

type ReviewDecision struct {
	RecordID string
	Accepted bool
	Reason   string
}
type ValidationIssue struct {
	Field    string
	Message  string
	Severity string
}
type BatchPolicy struct {
	AllowedDivisors   []int
	RequirePositive   bool
	RequireAllPassing bool
	MaxRecords        int
}

func DefaultPolicy() BatchPolicy {
	return BatchPolicy{AllowedDivisors: append([]int(nil), DefaultDivisors...), RequirePositive: true, RequireAllPassing: true, MaxRecords: 500}
}
func (p BatchPolicy) Validate() error {
	if len(p.AllowedDivisors) == 0 || p.MaxRecords <= 0 {
		return ErrInvalid
	}
	for _, divisor := range p.AllowedDivisors {
		if divisor < 2 {
			return fmt.Errorf("invalid divisor: %d", divisor)
		}
	}
	return nil
}
func (p BatchPolicy) CheckRecord(record StickerRecord) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if p.RequirePositive && record.Number <= 0 {
		issues = append(issues, ValidationIssue{Field: "number", Message: "number must be positive", Severity: "error"})
	}
	if record.ID == "" {
		issues = append(issues, ValidationIssue{Field: "id", Message: "id is required", Severity: "error"})
	}
	if record.Result != "pass" && record.Result != "fail" {
		issues = append(issues, ValidationIssue{Field: "result", Message: "result must be pass or fail", Severity: "warning"})
	}
	return issues
}
func (p BatchPolicy) CheckBatch(batch Batch) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if err := p.Validate(); err != nil {
		issues = append(issues, ValidationIssue{Field: "policy", Message: err.Error(), Severity: "error"})
	}
	if len(batch.Records) > p.MaxRecords {
		issues = append(issues, ValidationIssue{Field: "records", Message: "record limit exceeded", Severity: "error"})
	}
	for _, record := range batch.Records {
		issues = append(issues, p.CheckRecord(record)...)
	}
	return issues
}
func ApplyDecision(batch *Batch, decision ReviewDecision) error {
	if batch == nil || decision.RecordID == "" {
		return ErrInvalid
	}
	for index := range batch.Records {
		if batch.Records[index].ID != decision.RecordID {
			continue
		}
		batch.Records[index].Confirmed = decision.Accepted
		if decision.Reason != "" {
			batch.Records[index].UpdatedBy = decision.Reason
		}
		return nil
	}
	return ErrNotFound
}
func ReviewSummary(batch Batch) string {
	passing, failing, confirmed := CountResults(batch.Records)
	return fmt.Sprintf("%s: total=%d passing=%d failing=%d confirmed=%d", batch.Label, len(batch.Records), passing, failing, confirmed)
}
func NormalizeDivisors(divisors []int) []int {
	result := make([]int, 0, len(divisors))
	seen := map[int]bool{}
	for _, divisor := range divisors {
		if divisor > 1 && !seen[divisor] {
			seen[divisor] = true
			result = append(result, divisor)
		}
	}
	sort.Ints(result)
	return result
}
func ExplainResult(record StickerRecord) string {
	if record.Result == "pass" {
		return fmt.Sprintf("%d divisible by %s", record.Number, FormatDivisors(record.Divisors))
	}
	return fmt.Sprintf("%d has no configured divisor", record.Number)
}
func FormatDivisors(divisors []int) string {
	normalized := NormalizeDivisors(divisors)
	values := make([]string, len(normalized))
	for i, divisor := range normalized {
		values[i] = fmt.Sprintf("%d", divisor)
	}
	return strings.Join(values, ",")
}
func CloneRecords(records []StickerRecord) []StickerRecord {
	result := make([]StickerRecord, len(records))
	copy(result, records)
	for i := range result {
		result[i].Divisors = append([]int(nil), records[i].Divisors...)
	}
	return result
}
func SortRecords(records []StickerRecord) {
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
}
func CompareRecords(left, right StickerRecord) bool {
	return left.ID == right.ID && left.Number == right.Number && left.Result == right.Result && left.Confirmed == right.Confirmed
}
