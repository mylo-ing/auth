// api/controllers/mocks.go
package controllers

import (
	"sync"
	"time"

	"auth-service/api/infra/cache"
	"auth-service/api/infra/mailer"
)

// ----- memCache -----------------------------------------------------

type memCache struct {
	sync.RWMutex
	data map[string]entry
}
type entry struct {
	v   string
	exp time.Time
}

func newMemCache() *memCache { return &memCache{data: map[string]entry{}} }

func (m *memCache) SetValue(k, v string, ttl time.Duration) error {
	m.Lock()
	m.data[k] = entry{v, time.Now().Add(ttl)}
	m.Unlock()
	return nil
}
func (m *memCache) GetValue(k string) (string, error) {
	m.RLock()
	it, ok := m.data[k]
	m.RUnlock()
	if !ok || time.Now().After(it.exp) {
		return "", nil
	}
	return it.v, nil
}
func (m *memCache) DeleteKey(k string) error {
	m.Lock()
	delete(m.data, k)
	m.Unlock()
	return nil
}

// confirm interfaces
var _ cache.CodeCache = (*memCache)(nil)

// ----- noop mailer ---------------------------------------------------

type noopMailer struct{}

func (noopMailer) SendSignupConfirmation(string, string) error { return nil }

// confirm interface
var _ mailer.EmailSender = (*noopMailer)(nil)
