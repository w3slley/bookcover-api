package mocks

import (
	"bookcover-api/internal/cache"

	"github.com/bradfitz/gomemcache/memcache"
)

// MockMemcacheClient implements CacheClient interface for testing
type MockMemcacheClient struct {
	items map[string]*memcache.Item
}

// NewMockCache creates a new mock cache instance for testing
func NewMockCache() cache.CacheClient {
	return &MockMemcacheClient{
		items: make(map[string]*memcache.Item),
	}
}

func (m *MockMemcacheClient) Get(key string) (*memcache.Item, error) {
	if item, exists := m.items[key]; exists {
		return item, nil
	}
	return nil, memcache.ErrCacheMiss
}

func (m *MockMemcacheClient) Set(item *memcache.Item) error {
	m.items[item.Key] = item
	return nil
}

// Reset clears all items from the mock cache
func (m *MockMemcacheClient) Reset() {
	m.items = make(map[string]*memcache.Item)
}
