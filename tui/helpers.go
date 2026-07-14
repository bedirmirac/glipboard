package tui

func maxCharOfString(s string) string {
	if len(s) <= 25 {
		return s
	}
	runes := []rune(s)
	return string(runes[:25]) + "..."
}
