package markdown

import (
	"strings"
)

// Escape escapes special Markdown characters in a text string.
func Escape(s string) string {
	// Markdown characters that might trigger formatting.
	const specialChars = "\\`*_{}[]()#+-.!"
	var sb strings.Builder
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			sb.WriteRune('\\')
		}
		sb.WriteRune(r)
	}
	return sb.String()
}
