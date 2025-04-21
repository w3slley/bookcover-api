package routes

import (
	"bookcover-api/internal/routes"
	"bookcover-api/internal/cache"
	"bookcover-api/internal/helpers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

var isbn = "978-0345376597"
var expectedUrl = "https://example.com/book.jpg"

func setupTestCache() {
	// Initialize real Memcached client
	mc := memcache.New("localhost:11211")
	cache.SetCache(mc)
}

func cleanupTestCache() {
	// Clear the cache after tests
	if mc := cache.GetCache(); mc != nil {
		cache.SetCache(nil)
	}
}

func TestBookcoverSearch_CacheHit(t *testing.T) {
	setupTestCache()
	defer cleanupTestCache()

	// Setup test data
	cacheKey := "test+book+test+author"
	cache.GetCache().Set(&memcache.Item{Key: cacheKey, Value: []byte(expectedUrl)})

	// Create request
	req := httptest.NewRequest("GET", "/bookcover?book_title=test+book&author_name=test+author", nil)
	w := httptest.NewRecorder()

	// Call handler
	routes.BookcoverSearch(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["url"] != expectedUrl {
		t.Errorf("Expected URL %s, got %s", expectedUrl, response["url"])
	}
}

func TestBookcoverByIsbn_CacheHit(t *testing.T) {
	setupTestCache()
	defer cleanupTestCache()

	// Setup test data
	cacheKey := strings.ReplaceAll(isbn, "-", "")
	cache.GetCache().Set(&memcache.Item{Key: cacheKey, Value: []byte(expectedUrl)})

	// Create request
	req := httptest.NewRequest("GET", "/bookcover/"+isbn, nil)
	w := httptest.NewRecorder()

	// Call handler
	routes.BookcoverByIsbn(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["url"] != expectedUrl {
		t.Errorf("Expected URL %s, got %s", expectedUrl, response["url"])
	}
}

func TestBookcoverSearch_CacheMiss(t *testing.T) {
	setupTestCache()
	defer cleanupTestCache()

	// Create request
	req := httptest.NewRequest("GET", "/bookcover?book_title=test+book&author_name=test+author+different", nil)
	w := httptest.NewRecorder()

	// Call handler
	routes.BookcoverSearch(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 for cache miss, got %d", resp.StatusCode)
	}
}

func TestBookcoverByIsbn_CacheMiss(t *testing.T) {
	setupTestCache()
	defer cleanupTestCache()

	// Create request
	req := httptest.NewRequest("GET", "/bookcover/978-0000000000", nil)
	w := httptest.NewRecorder()

	// Call handler
	routes.BookcoverByIsbn(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
    t.Errorf("Expected status code 404 for cache miss, got %d", resp.StatusCode)
	}
	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 for cache miss, got %d", resp.StatusCode)
	}
}

func TestBookcoverSearch_InvalidParams(t *testing.T) {
	setupTestCache()
	defer cleanupTestCache()

	// Create request with missing parameters
	req := httptest.NewRequest("GET", "/bookcover", nil)
	w := httptest.NewRecorder()

	// Call handler
	routes.BookcoverSearch(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid params, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != helpers.MANDATORY_PARAMS_MISSING {
		t.Errorf("Expected error message %s, got %s", helpers.MANDATORY_PARAMS_MISSING, response["error"])
	}
}

func TestBookcoverByIsbn_InvalidISBN(t *testing.T) {
	setupTestCache()
	defer cleanupTestCache()

	// Create request with invalid ISBN
	req := httptest.NewRequest("GET", "/bookcover/123", nil)
	w := httptest.NewRecorder()

	// Call handler
	routes.BookcoverByIsbn(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid ISBN, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	if response["error"] != helpers.INVALID_ISBN {
		t.Errorf("Expected error message %s, got %s", helpers.INVALID_ISBN, response["error"])
	}
} 