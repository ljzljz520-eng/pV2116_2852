package importer

import (
	"strconv"
	"strings"
)

func WarningsText(result Result) string {
	values := append([]string(nil), result.Warnings...)
	return strings.Join(values, "; ")
}
func CandidateCount(result Result) int { return len(result.Candidates) }
func WarningCount(result Result) int   { return len(result.Warnings) }
func HasWarnings(result Result) bool   { return len(result.Warnings) > 0 }
func ResultLabel(result Result) string {
	if HasWarnings(result) {
		return "partial"
	}
	return "accepted"
}
func ResultCode(result Result) string {
	return strconv.Itoa(len(result.Candidates)) + "/" + strconv.Itoa(len(result.Warnings))
}
