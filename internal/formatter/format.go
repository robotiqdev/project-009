package formatter

import (
	"strconv"
	"strings"
)

// FormatResult formats a float64 value as a string with at most 2 decimal places,
// trimming trailing zeros and the decimal point for whole numbers.
func FormatResult(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
