package tui

func maxCharOfString(s string) string {
	runes := []rune(s)
	return string(runes[:25]) + "..."
}
