package importer

import (
	"bufio"
	"fmt"
	"io"
	"stickerchallenge/internal/domain"
	"strconv"
	"strings"
)

type Result struct {
	Candidates []domain.Candidate
	Warnings   []string
}

func ParseRows(input string) (Result, error) {
	return ParseReader(strings.NewReader(input))
}

func ParseReader(reader io.Reader) (Result, error) {
	result := Result{Candidates: make([]domain.Candidate, 0), Warnings: make([]string, 0)}
	scanner := bufio.NewScanner(reader)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		parts := strings.Split(text, ",")
		if len(parts) != 2 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("line %d: expected id,number", line))
			continue
		}
		number, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || number <= 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("line %d: invalid number", line))
			continue
		}
		id := strings.TrimSpace(parts[0])
		if id == "" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("line %d: missing id", line))
			continue
		}
		result.Candidates = append(result.Candidates, domain.Candidate{ID: id, Number: number})
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func ValidateCandidates(candidates []domain.Candidate) []string {
	warnings := make([]string, 0)
	seen := map[string]bool{}
	for _, candidate := range candidates {
		if seen[candidate.ID] {
			warnings = append(warnings, "duplicate id: "+candidate.ID)
		}
		seen[candidate.ID] = true
		if candidate.Number <= 0 {
			warnings = append(warnings, "non-positive number: "+candidate.ID)
		}
	}
	return warnings
}

func Normalize(candidates []domain.Candidate) []domain.Candidate {
	result := append([]domain.Candidate(nil), candidates...)
	sortCandidates(result)
	return result
}

func sortCandidates(candidates []domain.Candidate) {
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].ID < candidates[i].ID {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
}
