package cache

import (
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
)

var cache *memcache.Client

func GetCache() *memcache.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if cache != nil {
		return cache
	}
	cache = memcache.New(os.Getenv("MEMCACHED_HOST") + ":11211")
	log.Print("memcached connection established!")
	return cache
}
