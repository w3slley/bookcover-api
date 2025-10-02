package service

import (
	"errors"
	"testing"

	"bookcover-api/internal/scraper"
	"bookcover-api/mocks"

	"github.com/bradfitz/gomemcache/memcache"
)

type mockScraper struct {
	fetchByTitleAuthorFunc func(bookTitle, authorName string) (string, error)
	fetchByISBNFunc        func(isbn string) (string, error)
}

func (m *mockScraper) FetchByTitleAuthor(bookTitle, authorName string) (string, error) {
	if m.fetchByTitleAuthorFunc != nil {
		return m.fetchByTitleAuthorFunc(bookTitle, authorName)
	}
	return "", errors.New("not implemented")
}

func (m *mockScraper) FetchByISBN(isbn string) (string, error) {
	if m.fetchByISBNFunc != nil {
		return m.fetchByISBNFunc(isbn)
	}
	return "", errors.New("not implemented")
}

// Ensure mockScraper implements scraper.Scraper interface
var _ scraper.Scraper = (*mockScraper)(nil)

func TestNewBookcoverService(t *testing.T) {
	mockScraper := &mockScraper{}
	mockCache := mocks.NewMockCache()

	service := NewBookcoverService(mockScraper, mockCache)

	if service == nil {
		t.Error("NewBookcoverService() returned nil")
	}
}

func TestGetByTitleAuthor_CacheHit(t *testing.T) {
	mockScraper := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			t.Error("Scraper should not be called when cache hit occurs")
			return "", nil
		},
	}

	mockCache := mocks.NewMockCache()
	expectedURL := "https://example.com/cached-cover.jpg"
	mockCache.Set(&memcache.Item{
		Key:   "test+book+test+author",
		Value: []byte(expectedURL),
	})

	service := NewBookcoverService(mockScraper, mockCache)

	url, err := service.GetByTitleAuthor("test book", "test author")
	if err != nil {
		t.Errorf("GetByTitleAuthor() error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expectedURL)
	}
}

func TestGetByTitleAuthor_CacheMiss_ScraperSuccess(t *testing.T) {
	expectedURL := "https://example.com/scraped-cover.jpg"

	mockScraper := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			if bookTitle != "test+book" || authorName != "test+author" {
				t.Errorf("Scraper received wrong parameters: title=%v, author=%v", bookTitle, authorName)
			}
			return expectedURL, nil
		},
	}

	mockCache := mocks.NewMockCache()
	service := NewBookcoverService(mockScraper, mockCache)

	url, err := service.GetByTitleAuthor("test book", "test author")
	if err != nil {
		t.Errorf("GetByTitleAuthor() error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expectedURL)
	}

	cachedItem, _ := mockCache.Get("test+book+test+author")
	if cachedItem == nil {
		t.Error("Expected item to be cached, but cache is empty")
	} else if string(cachedItem.Value) != expectedURL {
		t.Errorf("Cached value = %v, want %v", string(cachedItem.Value), expectedURL)
	}
}

func TestGetByTitleAuthor_ScraperError(t *testing.T) {
	expectedError := errors.New("scraper failed")

	mockScraper := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			return "", expectedError
		},
	}

	mockCache := mocks.NewMockCache()
	service := NewBookcoverService(mockScraper, mockCache)

	_, err := service.GetByTitleAuthor("test book", "test author")
	if err == nil {
		t.Error("GetByTitleAuthor() expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("GetByTitleAuthor() error = %v, want %v", err, expectedError)
	}
}

func TestGetByISBN_CacheHit(t *testing.T) {
	mockScraper := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			t.Error("Scraper should not be called when cache hit occurs")
			return "", nil
		},
	}

	mockCache := mocks.NewMockCache()
	expectedURL := "https://example.com/cached-isbn-cover.jpg"
	mockCache.Set(&memcache.Item{
		Key:   "9780345376596",
		Value: []byte(expectedURL),
	})

	service := NewBookcoverService(mockScraper, mockCache)

	url, err := service.GetByISBN("978-0345376596")
	if err != nil {
		t.Errorf("GetByISBN() error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByISBN() = %v, want %v", url, expectedURL)
	}
}

func TestGetByISBN_CacheMiss_ScraperSuccess(t *testing.T) {
	expectedURL := "https://example.com/isbn-cover.jpg"

	mockScraper := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			if isbn != "9780345376596" {
				t.Errorf("Scraper received wrong ISBN: %v", isbn)
			}
			return expectedURL, nil
		},
	}

	mockCache := mocks.NewMockCache()
	service := NewBookcoverService(mockScraper, mockCache)

	url, err := service.GetByISBN("978-0345376596")
	if err != nil {
		t.Errorf("GetByISBN() error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByISBN() = %v, want %v", url, expectedURL)
	}

	// Verify cache was set (ISBN is normalized to remove dashes and lowercase)
	cachedItem, _ := mockCache.Get("9780345376596")
	if cachedItem == nil {
		t.Error("Expected item to be cached, but cache is empty")
	} else if string(cachedItem.Value) != expectedURL {
		t.Errorf("Cached value = %v, want %v", string(cachedItem.Value), expectedURL)
	}
}

func TestGetByISBN_ScraperError(t *testing.T) {
	expectedError := errors.New("ISBN not found")

	mockScraper := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			return "", expectedError
		},
	}

	mockCache := mocks.NewMockCache()
	service := NewBookcoverService(mockScraper, mockCache)

	_, err := service.GetByISBN("978-0000000000")
	if err == nil {
		t.Error("GetByISBN() expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("GetByISBN() error = %v, want %v", err, expectedError)
	}
}

func TestGetByTitleAuthor_NilCache(t *testing.T) {
	expectedURL := "https://example.com/no-cache-cover.jpg"

	mockScraper := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			return expectedURL, nil
		},
	}

	service := NewBookcoverService(mockScraper, nil)

	url, err := service.GetByTitleAuthor("test book", "test author")
	if err != nil {
		t.Errorf("GetByTitleAuthor() with nil cache error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expectedURL)
	}
}

func TestGetByISBN_NilCache(t *testing.T) {
	expectedURL := "https://example.com/no-cache-isbn.jpg"

	mockScraper := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			return expectedURL, nil
		},
	}

	service := NewBookcoverService(mockScraper, nil)

	url, err := service.GetByISBN("978-0345376596")
	if err != nil {
		t.Errorf("GetByISBN() with nil cache error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByISBN() = %v, want %v", url, expectedURL)
	}
}
