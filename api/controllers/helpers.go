package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type SubscriberResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func validateEmail(email string) error {
	if email == "" || !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid or missing email")
	}
	return nil
}

// Generate a random 6-digit numeric code
func generateSixDigitCode() string {
	var b [3]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Println("Failed to generate random bytes, fallback to time-based code.")
		now := time.Now().UnixNano()
		return fmt.Sprintf("%06d", now%1000000)
	}
	num := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", num)
}

// randomToken returns a URL-safe random string
func randomToken(length int) string {
	raw := make([]byte, length)
	_, _ = rand.Read(raw)
	return base64.RawURLEncoding.EncodeToString(raw)
}
