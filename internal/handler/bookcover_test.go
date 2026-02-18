package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bookcover-api/internal/cache"
	"bookcover-api/internal/config"
	"bookcover-api/internal/scraper"
	"bookcover-api/internal/service"
	"bookcover-api/mocks"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	isbn        = "978-0345376597"
	expectedURL = "https://example.com/book.jpg"
)

func setupTestHandler() (*BookcoverHandler, cache.CacheClient) {
	mockCache := mocks.NewMockCache()
	goodreadsScraper := scraper.NewGoodreads()
	bookcoverService := service.NewBookcoverService(goodreadsScraper, mockCache)
	handler := NewBookcoverHandler(bookcoverService)
	return handler, mockCache
}

func TestBookcoverSearch_CacheHit(t *testing.T) {
	handler, mockCache := setupTestHandler()

	// Setup test data
	cacheKey := "test+book+test+author"
	mockCache.Set(&memcache.Item{Key: cacheKey, Value: []byte(expectedURL)})

	// Create request
	req := httptest.NewRequest("GET", "/bookcover?book_title=test+book&author_name=test+author", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Search(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["url"] != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, response["url"])
	}
}

func TestBookcoverByISBN_CacheHit(t *testing.T) {
	handler, mockCache := setupTestHandler()

	// Setup test data
	cacheKey := strings.ReplaceAll(isbn, "-", "")
	mockCache.Set(&memcache.Item{Key: cacheKey, Value: []byte(expectedURL)})

	// Create request
	req := httptest.NewRequest("GET", "/bookcover/"+isbn, nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ByISBN(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["url"] != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, response["url"])
	}
}

func TestBookcoverSearch_CacheMiss(t *testing.T) {
	handler, _ := setupTestHandler()

	// Create request
	req := httptest.NewRequest("GET", "/bookcover?book_title=test+book&author_name=test+author+different", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Search(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 for cache miss, got %d", resp.StatusCode)
	}
}

func TestBookcoverByISBN_CacheMiss(t *testing.T) {
	handler, _ := setupTestHandler()

	// Create request
	req := httptest.NewRequest("GET", "/bookcover/978-0000000000", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ByISBN(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 for cache miss, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
}

func TestBookcoverSearch_ByISBNQueryParam_CacheHit(t *testing.T) {
	handler, mockCache := setupTestHandler()

	cacheKey := strings.ReplaceAll(isbn, "-", "")
	mockCache.Set(&memcache.Item{Key: cacheKey, Value: []byte(expectedURL)})

	req := httptest.NewRequest("GET", "/bookcover?isbn="+isbn, nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["url"] != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, response["url"])
	}
}

func TestBookcoverSearch_ByISBNQueryParam_InvalidISBN(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/bookcover?isbn=123", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid ISBN, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.InvalidISBN {
		t.Errorf("Expected error message %s, got %s", config.InvalidISBN, response["error"])
	}
}

func TestBookcoverSearch_ConflictingParams(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/bookcover?isbn=978-0345376597&book_title=test", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for conflicting params, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.ConflictingParams {
		t.Errorf("Expected error message %s, got %s", config.ConflictingParams, response["error"])
	}
}

func TestBookcoverSearch_InvalidParams(t *testing.T) {
	handler, _ := setupTestHandler()

	// Create request with missing parameters
	req := httptest.NewRequest("GET", "/bookcover", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Search(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid params, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.MandidatoryParamsMissing {
		t.Errorf("Expected error message %s, got %s", config.MandidatoryParamsMissing, response["error"])
	}
}

func TestBookcoverByISBN_InvalidISBN(t *testing.T) {
	handler, _ := setupTestHandler()

	// Create request with invalid ISBN
	req := httptest.NewRequest("GET", "/bookcover/123", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ByISBN(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid ISBN, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.InvalidISBN {
		t.Errorf("Expected error message %s, got %s", config.InvalidISBN, response["error"])
	}
}

func TestBookcoverSearch_OnlyBookTitle_MissingAuthor(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/bookcover?book_title=test+book", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 when author_name is missing, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.MandidatoryParamsMissing {
		t.Errorf("Expected error %s, got %s", config.MandidatoryParamsMissing, response["error"])
	}
}

func TestBookcoverSearch_OnlyAuthorName_MissingBookTitle(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/bookcover?author_name=test+author", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 when book_title is missing, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.MandidatoryParamsMissing {
		t.Errorf("Expected error %s, got %s", config.MandidatoryParamsMissing, response["error"])
	}
}

func TestBookcoverSearch_ISBNWithAuthorName_ConflictingParams(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/bookcover?isbn=978-0345376597&author_name=test+author", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 for isbn+author_name conflict, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != config.ConflictingParams {
		t.Errorf("Expected error %s, got %s", config.ConflictingParams, response["error"])
	}
}
