package domain

import (
	"fmt"
	"sort"
)

var DefaultDivisors = []int{2, 3, 5, 7, 11}

func EvaluateNumber(number int, divisors []int) (string, []int, error) {
	if number <= 0 {
		return "invalid", nil, fmt.Errorf("number %d must be positive: %w", number, ErrInvalid)
	}
	if len(divisors) == 0 {
		return "invalid", nil, fmt.Errorf("divisors are required: %w", ErrInvalid)
	}
	hits := make([]int, 0, len(divisors))
	seen := map[int]bool{}
	for _, divisor := range divisors {
		if divisor <= 1 || seen[divisor] {
			continue
		}
		seen[divisor] = true
		if number%divisor == 0 {
			hits = append(hits, divisor)
		}
	}
	sort.Ints(hits)
	if len(hits) == 0 {
		return "fail", hits, nil
	}
	return "pass", hits, nil
}

func Recalculate(record StickerRecord, divisors []int) (StickerRecord, error) {
	result, hits, err := EvaluateNumber(record.Number, divisors)
	if err != nil {
		return record, err
	}
	record.Result = result
	record.Divisors = hits
	record.Confirmed = false
	return record, nil
}

func ValidateBatch(batch Batch) error {
	if batch.ID == "" || batch.Label == "" || batch.Owner == "" {
		return ErrInvalid
	}
	if len(batch.Records) == 0 {
		return fmt.Errorf("batch has no records: %w", ErrInvalid)
	}
	for _, record := range batch.Records {
		if record.BatchID != batch.ID || record.ID == "" || record.Number <= 0 {
			return ErrInvalid
		}
	}
	return nil
}

func CountResults(records []StickerRecord) (int, int, int) {
	passing, failing, confirmed := 0, 0, 0
	for _, record := range records {
		if record.Result == "pass" {
			passing++
		} else if record.Result == "fail" {
			failing++
		}
		if record.Confirmed {
			confirmed++
		}
	}
	return passing, failing, confirmed
}
