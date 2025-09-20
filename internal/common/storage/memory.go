package storage

import "sync"

// InMemory is a threadsafe map-backed storage helper.
type InMemory[T any] struct {
	mu sync.RWMutex
	db map[string]T
}

// NewInMemory creates a new store.
func NewInMemory[T any]() *InMemory[T] {
	return &InMemory[T]{db: make(map[string]T)}
}

// Get returns a value by key.
func (m *InMemory[T]) Get(key string) (T, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.db[key]
	return val, ok
}

// Set stores a value for key.
func (m *InMemory[T]) Set(key string, value T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db[key] = value
}

// Delete removes the value for key.
func (m *InMemory[T]) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.db, key)
}

// List returns a snapshot of the current values.
func (m *InMemory[T]) List() []T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]T, 0, len(m.db))
	for _, value := range m.db {
		out = append(out, value)
	}
	return out
}
