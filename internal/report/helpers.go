package report

import (
	"strconv"
	"strings"
)

func FormatCount(value int) string       { return strconv.Itoa(value) }
func IsEmpty(value string) bool          { return value == "" }
func NormalizeLabel(value string) string { return strings.TrimSpace(value) }
