package mocks

import (
	"fmt"

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

func (m *MockMemcacheClient) Add(item *memcache.Item) error {
	if _, exists := m.items[item.Key]; exists {
		return memcache.ErrNotStored
	}
	m.items[item.Key] = item
	return nil
}

func (m *MockMemcacheClient) Increment(key string, delta uint64) (uint64, error) {
	item, exists := m.items[key]
	if !exists {
		return 0, memcache.ErrCacheMiss
	}
	val := uint64(0)
	for _, b := range item.Value {
		val = val*10 + uint64(b-'0')
	}
	val += delta
	item.Value = fmt.Appendf(nil, "%d", val)
	return val, nil
}

// Reset clears all items from the mock cache
func (m *MockMemcacheClient) Reset() {
	m.items = make(map[string]*memcache.Item)
}
