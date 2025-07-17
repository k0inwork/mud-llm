package testutils

import (
	"time"
)

// MockCache implements the CacheInterface for testing purposes.
type MockCache struct {
	Data map[string]interface{}
}

func NewMockCache() *MockCache {
	return &MockCache{Data: make(map[string]interface{})}
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	val, found := m.Data[key]
	return val, found
}

func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) {
	m.Data[key] = value
}

func (m *MockCache) SetMany(data map[string]interface{}, ttl time.Duration) {
	for k, v := range data {
		m.Data[k] = v
	}
}

func (m *MockCache) Delete(key string) {
	delete(m.Data, key)
}

func (m *MockCache) Clear() {
	m.Data = make(map[string]interface{})
}
