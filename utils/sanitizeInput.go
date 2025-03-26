package utils

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeInput removes HTML tags, scripts, and trims spaces
func SanitizeInput(input string) string {
	// Decode HTML entities to prevent encoded attacks
	input = html.UnescapeString(input)

	// Remove script tags and content inside
	re := regexp.MustCompile(`(?i)<script.*?>.*?</script>`)
	input = re.ReplaceAllString(input, "")

	// Remove all other HTML tags
	re = regexp.MustCompile(`(?i)<.*?>`)
	input = re.ReplaceAllString(input, "")

	// Trim extra spaces
	input = strings.TrimSpace(input)

	return input
}
