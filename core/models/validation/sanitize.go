package Sanitize

import (
    "github.com/microcosm-cc/bluemonday"
	"regexp"
)

func SanitizeInput(input string) string {
    // Use bluemonday to sanitize HTML content
    p := bluemonday.UGCPolicy()
    return p.Sanitize(input)
}

func ValidateMessage(message string) bool {
    // Example regex to allow only printable characters and spaces
    re := regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{P}\p{Zs}]+$`)
    return re.MatchString(message)
}