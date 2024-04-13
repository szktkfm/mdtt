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
