package cache

import (
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

type CacheClient interface {
	Get(key string) (*memcache.Item, error)
	Set(item *memcache.Item) error
}

var cache CacheClient

func GetCache() CacheClient {
	if cache != nil {
		return cache
	}
	cache = memcache.New(os.Getenv("MEMCACHED_HOST") + ":11211")
	return cache
}

func SetCache(c CacheClient) {
	cache = c
}
