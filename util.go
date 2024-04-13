package mdtt

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func padOrTruncate(s string, n int) string {
	if runewidth.StringWidth(s) > n {
		return runewidth.Truncate(s, n, "")
	} else {
		return s + strings.Repeat(" ", n-runewidth.StringWidth(s))
	}
}
func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
