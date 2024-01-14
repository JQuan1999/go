package ch1

import (
	"strings"
)

func Count(s string, substr string) int {
	count := strings.Count(s, substr)
	return count
}
