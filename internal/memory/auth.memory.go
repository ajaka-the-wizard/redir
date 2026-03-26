package memory

import (
	"sync"

	"github.com/ajaka-the-wizard/redir/internal/domain"
)

type AuthMemoryMap struct {
	mu   sync.RWMutex
	auth map[string]domain.
		LightUser
}

func NewMemoryMap() *AuthMemoryMap {
	return &AuthMemoryMap{
		auth: map[string]domain.
			LightUser{},
	}
}

func (m *AuthMemoryMap) SetUserOnline(sessionId string, u *domain.
	LightUser) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auth[sessionId] = *u
}

func (m *AuthMemoryMap) GetUser(sessionId string) (*domain.
	LightUser, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if u, ok := m.auth[sessionId]; ok {
		return &u, true
	}
	return nil, false
}

func (m *AuthMemoryMap) RevokeUser(sessionId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.auth, sessionId)
}
