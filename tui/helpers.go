package tui

import "strings"

func maxCharOfString(s string) string {
	first, _, ok := strings.Cut(s, "\n")
	if len(first) <= 25 && !ok {
		return first
	}
	runes := []rune(first)
	return string(runes[:25]) + "..."
}
