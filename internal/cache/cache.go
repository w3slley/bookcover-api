package cache

import (
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

var cache *memcache.Client

func GetCache() *memcache.Client {
	if cache != nil {
		return cache
	}
	cache = memcache.New(os.Getenv("MEMCACHED_HOST") + ":11211")
	return cache
}
