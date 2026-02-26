package service

import (
	"log"
	"strings"

	"bookcover-api/internal/cache"
	"bookcover-api/internal/scraper"

	"github.com/bradfitz/gomemcache/memcache"
)

const querySeparator = "+"

type bookcoverService struct {
	scraper scraper.Scraper
	cache   cache.CacheClient
}

func NewBookcoverService(s scraper.Scraper, cache cache.CacheClient) BookcoverService {
	return &bookcoverService{
		scraper: s,
		cache:   cache,
	}
}

func (s *bookcoverService) GetByTitleAuthor(bookTitle, authorName, imageSize string) (string, error) {
	bookTitle = strings.ReplaceAll(bookTitle, " ", querySeparator)
	authorName = strings.ReplaceAll(authorName, " ", querySeparator)
	cacheKey := strings.ToLower(bookTitle + querySeparator + authorName)

	if cachedURL, err := s.getFromCache(cacheKey); cachedURL != "" {
		return applyImageSize(cachedURL, imageSize), err
	}

	imageURL, err := s.scraper.FetchByTitleAuthor(bookTitle, authorName)
	if err != nil {
		return "", err
	}

	s.setCache(cacheKey, imageURL)

	return applyImageSize(imageURL, imageSize), nil
}

func (s *bookcoverService) GetByISBN(isbn, imageSize string) (string, error) {
	isbn = strings.ReplaceAll(isbn, "-", "")
	cacheKey := strings.ToLower(isbn)

	if cachedURL, err := s.getFromCache(cacheKey); cachedURL != "" {
		return applyImageSize(cachedURL, imageSize), err
	}

	imageURL, err := s.scraper.FetchByISBN(isbn)
	if err != nil {
		return "", err
	}

	s.setCache(cacheKey, imageURL)

	return applyImageSize(imageURL, imageSize), nil
}

func applyImageSize(url, imageSize string) string {
	switch imageSize {
	case "small":
		return insertSizeSuffix(url, "__SY75__")
	case "medium":
		return insertSizeSuffix(url, "__SY375__")
	default:
		return url
	}
}

func insertSizeSuffix(url, suffix string) string {
	dotIndex := strings.LastIndex(url, ".")
	if dotIndex == -1 {
		return url
	}
	return url[:dotIndex] + "." + suffix + url[dotIndex:]
}

func (s *bookcoverService) getFromCache(key string) (string, error) {
	if s.cache == nil {
		return "", nil
	}

	cachedURL, err := s.cache.Get(key)
	if err != nil {
		log.Print(err)
		return "", nil
	}

	if cachedURL != nil {
		log.Printf("Found cache with key %s", key)
		return string(cachedURL.Value), nil
	}

	return "", nil
}

func (s *bookcoverService) setCache(key, value string) {
	if s.cache == nil {
		return
	}

	err := s.cache.Set(&memcache.Item{Key: key, Value: []byte(value)})
	if err != nil {
		log.Printf("Failed to set cache for key %s: %v", key, err)
		return
	}

	log.Printf("Created cache for key %s", key)
}
