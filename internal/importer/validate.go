package importer

import (
	"fmt"
	"sort"
	"stickerchallenge/internal/domain"
)

type ValidationReport struct {
	Total    int
	Valid    int
	Invalid  int
	Warnings []string
	IDs      []string
}

func BuildReport(result Result) ValidationReport {
	warnings := append([]string(nil), result.Warnings...)
	report := ValidationReport{Total: len(result.Candidates) + len(warnings), Valid: len(result.Candidates), Invalid: len(warnings), Warnings: warnings}
	for _, candidate := range result.Candidates {
		report.IDs = append(report.IDs, candidate.ID)
	}
	sort.Strings(report.IDs)
	return report
}
func ValidateAgainstPolicy(candidates []domain.Candidate, policy domain.BatchPolicy) ValidationReport {
	report := ValidationReport{Total: len(candidates), Warnings: make([]string, 0)}
	seen := map[string]bool{}
	for _, candidate := range candidates {
		invalid := false
		if candidate.ID == "" {
			report.Warnings = append(report.Warnings, "missing id")
			invalid = true
		}
		if candidate.Number <= 0 && policy.RequirePositive {
			report.Warnings = append(report.Warnings, fmt.Sprintf("%s has non-positive number", candidate.ID))
			invalid = true
		}
		if seen[candidate.ID] {
			report.Warnings = append(report.Warnings, fmt.Sprintf("duplicate id %s", candidate.ID))
			invalid = true
		}
		seen[candidate.ID] = true
		if invalid {
			report.Invalid++
		} else {
			report.Valid++
			report.IDs = append(report.IDs, candidate.ID)
		}
	}
	sort.Strings(report.IDs)
	return report
}
func SplitValid(candidates []domain.Candidate) ([]domain.Candidate, []domain.Candidate) {
	valid := make([]domain.Candidate, 0)
	invalid := make([]domain.Candidate, 0)
	for _, candidate := range candidates {
		if candidate.ID != "" && candidate.Number > 0 {
			valid = append(valid, candidate)
		} else {
			invalid = append(invalid, candidate)
		}
	}
	return valid, invalid
}
func Dedupe(candidates []domain.Candidate) []domain.Candidate {
	result := make([]domain.Candidate, 0, len(candidates))
	seen := map[string]bool{}
	for _, candidate := range candidates {
		if !seen[candidate.ID] {
			seen[candidate.ID] = true
			result = append(result, candidate)
		}
	}
	return result
}
func FindCandidate(candidates []domain.Candidate, id string) (domain.Candidate, bool) {
	for _, candidate := range candidates {
		if candidate.ID == id {
			return candidate, true
		}
	}
	return domain.Candidate{}, false
}
func ReplaceCandidate(candidates []domain.Candidate, replacement domain.Candidate) ([]domain.Candidate, bool) {
	result := append([]domain.Candidate(nil), candidates...)
	for index := range result {
		if result[index].ID == replacement.ID {
			result[index] = replacement
			return result, true
		}
	}
	return result, false
}
func CandidateNumbers(candidates []domain.Candidate) []int {
	result := make([]int, 0, len(candidates))
	for _, candidate := range candidates {
		result = append(result, candidate.Number)
	}
	sort.Ints(result)
	return result
}
