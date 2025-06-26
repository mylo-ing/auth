package cache

import "time"

// CodeCache is the minimal contract controllers need.
type CodeCache interface {
	SetValue(key, val string, ttl time.Duration) error
	GetValue(key string) (string, error)
	DeleteKey(key string) error
}
