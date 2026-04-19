package memory

import (
	"sync"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/domain"
)

type AuthMemoryMap struct {
	mu   sync.RWMutex
	auth map[string]*domain.
		LightUser
}

func NewMemoryMap() *AuthMemoryMap {
	return &AuthMemoryMap{
		auth: map[string]*domain.
			LightUser{},
	}
}

func (m *AuthMemoryMap) SetUserOnline(sessionId string, u *domain.
	LightUser) time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	u.LastAccessedTime = now
	u.Expires = now.Add(24 * time.Hour)
	m.auth[sessionId] = u
	return now
}

func (m *AuthMemoryMap) GetUser(sessionId string) (*domain.
	LightUser, bool) {
	m.mu.RLock()
	if u, ok := m.auth[sessionId]; ok {
		if time.Now().After(u.Expires) {
			m.mu.RUnlock()
			m.RevokeUser(sessionId)
			return nil, false
		}
		return u, true
	}
	m.mu.RUnlock()
	return nil, false
}

func (m *AuthMemoryMap) RevokeUser(sessionId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.auth, sessionId)
}

func (m *AuthMemoryMap) UpdateUserTimestamp(sessionId string) time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, ok := m.auth[sessionId]
	if !ok {
		return time.Time{}
	}
	now := time.Now()
	u.LastAccessedTime = now
	return now
}
