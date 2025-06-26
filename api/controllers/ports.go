// api/controllers/ports.go
package controllers

import "time"

type CodeCache interface {
	SetValue(key, val string, ttl time.Duration) error
	GetValue(key string) (string, error)
	DeleteKey(key string) error
}

type Mailer interface {
	SendSignupConfirmation(email, code string) error
}
