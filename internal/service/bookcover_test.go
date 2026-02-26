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

	url, err := service.GetByTitleAuthor("test book", "test author", "")
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

	url, err := service.GetByTitleAuthor("test book", "test author", "")
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

	_, err := service.GetByTitleAuthor("test book", "test author", "")
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

	url, err := service.GetByISBN("978-0345376596", "")
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

	url, err := service.GetByISBN("978-0345376596", "")
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

	_, err := service.GetByISBN("978-0000000000", "")
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

	url, err := service.GetByTitleAuthor("test book", "test author", "")
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

	url, err := service.GetByISBN("978-0345376596", "")
	if err != nil {
		t.Errorf("GetByISBN() with nil cache error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("GetByISBN() = %v, want %v", url, expectedURL)
	}
}

// errCache returns a non-ErrCacheMiss error on Get and an error on Set,
// allowing us to test the error-handling paths in getFromCache / setCache.
type errCache struct {
	getErr error
	setErr error
}

func (e *errCache) Get(key string) (*memcache.Item, error) {
	return nil, e.getErr
}

func (e *errCache) Set(item *memcache.Item) error {
	return e.setErr
}

func (e *errCache) Add(item *memcache.Item) error {
	return e.setErr
}

func (e *errCache) Increment(key string, delta uint64) (uint64, error) {
	return 0, e.getErr
}

func TestGetByTitleAuthor_CacheGetError_FallsBackToScraper(t *testing.T) {
	expectedURL := "https://example.com/scraped.jpg"

	ms := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			return expectedURL, nil
		},
	}

	svc := NewBookcoverService(ms, &errCache{getErr: errors.New("connection refused")})

	url, err := svc.GetByTitleAuthor("test book", "test author", "")
	if err != nil {
		t.Errorf("GetByTitleAuthor() unexpected error: %v", err)
	}
	if url != expectedURL {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expectedURL)
	}
}

func TestGetByISBN_CacheGetError_FallsBackToScraper(t *testing.T) {
	expectedURL := "https://example.com/isbn-scraped.jpg"

	ms := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			return expectedURL, nil
		},
	}

	svc := NewBookcoverService(ms, &errCache{getErr: errors.New("connection refused")})

	url, err := svc.GetByISBN("978-0345376596", "")
	if err != nil {
		t.Errorf("GetByISBN() unexpected error: %v", err)
	}
	if url != expectedURL {
		t.Errorf("GetByISBN() = %v, want %v", url, expectedURL)
	}
}

func TestGetByTitleAuthor_CacheSetError_StillReturnsURL(t *testing.T) {
	expectedURL := "https://example.com/scraped.jpg"

	ms := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			return expectedURL, nil
		},
	}

	// Get succeeds with a miss (nil, ErrCacheMiss would be handled; here we use
	// an error on Get to skip cache, and an error on Set to exercise that path).
	svc := NewBookcoverService(ms, &errCache{
		getErr: errors.New("get error"),
		setErr: errors.New("set error"),
	})

	url, err := svc.GetByTitleAuthor("test book", "test author", "")
	if err != nil {
		t.Errorf("GetByTitleAuthor() unexpected error: %v", err)
	}
	if url != expectedURL {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expectedURL)
	}
}

func TestGetByISBN_CacheSetError_StillReturnsURL(t *testing.T) {
	expectedURL := "https://example.com/isbn-scraped.jpg"

	ms := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			return expectedURL, nil
		},
	}

	svc := NewBookcoverService(ms, &errCache{
		getErr: errors.New("get error"),
		setErr: errors.New("set error"),
	})

	url, err := svc.GetByISBN("978-0345376596", "")
	if err != nil {
		t.Errorf("GetByISBN() unexpected error: %v", err)
	}
	if url != expectedURL {
		t.Errorf("GetByISBN() = %v, want %v", url, expectedURL)
	}
}

func TestApplyImageSize(t *testing.T) {
	baseURL := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.jpg"

	tests := []struct {
		name      string
		imageSize string
		expected  string
	}{
		{
			name:      "small image size",
			imageSize: "small",
			expected:  "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.__SY75__.jpg",
		},
		{
			name:      "medium image size",
			imageSize: "medium",
			expected:  "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.__SY375__.jpg",
		},
		{
			name:      "large image size returns original",
			imageSize: "large",
			expected:  baseURL,
		},
		{
			name:      "empty image size returns original",
			imageSize: "",
			expected:  baseURL,
		},
		{
			name:      "invalid image size returns original",
			imageSize: "xlarge",
			expected:  baseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyImageSize(baseURL, tt.imageSize)
			if result != tt.expected {
				t.Errorf("applyImageSize(%q, %q) = %q, want %q", baseURL, tt.imageSize, result, tt.expected)
			}
		})
	}
}

func TestApplyImageSize_PngExtension(t *testing.T) {
	url := "https://example.com/image.png"
	result := applyImageSize(url, "small")
	expected := "https://example.com/image.__SY75__.png"
	if result != expected {
		t.Errorf("applyImageSize() = %q, want %q", result, expected)
	}
}

func TestGetByTitleAuthor_WithImageSize(t *testing.T) {
	scraperURL := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.jpg"

	ms := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			return scraperURL, nil
		},
	}

	mockCache := mocks.NewMockCache()
	svc := NewBookcoverService(ms, mockCache)

	url, err := svc.GetByTitleAuthor("test book", "test author", "small")
	if err != nil {
		t.Errorf("GetByTitleAuthor() unexpected error: %v", err)
	}

	expected := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.__SY75__.jpg"
	if url != expected {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expected)
	}
}

func TestGetByISBN_WithImageSize(t *testing.T) {
	scraperURL := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.jpg"

	ms := &mockScraper{
		fetchByISBNFunc: func(isbn string) (string, error) {
			return scraperURL, nil
		},
	}

	mockCache := mocks.NewMockCache()
	svc := NewBookcoverService(ms, mockCache)

	url, err := svc.GetByISBN("978-0345376596", "medium")
	if err != nil {
		t.Errorf("GetByISBN() unexpected error: %v", err)
	}

	expected := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.__SY375__.jpg"
	if url != expected {
		t.Errorf("GetByISBN() = %v, want %v", url, expected)
	}
}

func TestGetByTitleAuthor_CacheHit_WithImageSize(t *testing.T) {
	cachedURL := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.jpg"

	ms := &mockScraper{
		fetchByTitleAuthorFunc: func(bookTitle, authorName string) (string, error) {
			t.Error("Scraper should not be called when cache hit occurs")
			return "", nil
		},
	}

	mockCache := mocks.NewMockCache()
	mockCache.Set(&memcache.Item{
		Key:   "test+book+test+author",
		Value: []byte(cachedURL),
	})

	svc := NewBookcoverService(ms, mockCache)

	url, err := svc.GetByTitleAuthor("test book", "test author", "small")
	if err != nil {
		t.Errorf("GetByTitleAuthor() unexpected error: %v", err)
	}

	expected := "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1555447414i/44767458.__SY75__.jpg"
	if url != expected {
		t.Errorf("GetByTitleAuthor() = %v, want %v", url, expected)
	}
}
