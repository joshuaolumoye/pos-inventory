package utils

import (
	"regexp"

	"math/rand"
	"time"

	"github.com/google/uuid"
)

// RandomDigits returns a string of n random digits
func RandomDigits(n int) string {
	rand.Seed(time.Now().UnixNano())
	digits := make([]byte, n)
	for i := 0; i < n; i++ {
		digits[i] = byte('0' + rand.Intn(10))
	}
	return string(digits)
}

// Sanitize strips out common XSS vectors (very basic, for demo; use a library for production)
func Sanitize(input string) string {
	// Remove script tags and angle brackets
	re := regexp.MustCompile(`(?i)<.*?script.*?>|<.*?>`)
	return re.ReplaceAllString(input, "")
}

// GenerateUUID returns a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}
