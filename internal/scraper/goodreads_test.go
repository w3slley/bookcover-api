package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const coverWithSizeSuffix = `[{"imageUrl": "https://example.com/cover._SX98_.jpg"}, {"imageUrl": "https://example.com/other.jpg"}]`

func newTestGoodreads(t *testing.T, handler http.HandlerFunc) *Goodreads {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &Goodreads{autoCompleteURL: srv.URL + "/?q="}
}

func respondWith(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}
}

func TestNewGoodreads(t *testing.T) {
	scraper := NewGoodreads()
	if scraper == nil {
		t.Error("NewGoodreads() returned nil")
	}
}

func TestFetchByISBN(t *testing.T) {
	var gotQuery string
	g := newTestGoodreads(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("q")
		w.Write([]byte(coverWithSizeSuffix))
	})

	url, err := g.FetchByISBN("9780345376596")
	if err != nil {
		t.Fatalf("FetchByISBN() error = %v", err)
	}

	if gotQuery != "9780345376596" {
		t.Errorf("FetchByISBN() queried q = %q, want %q", gotQuery, "9780345376596")
	}

	expectedURL := "https://example.com/cover.jpg"
	if url != expectedURL {
		t.Errorf("FetchByISBN() = %v, want %v", url, expectedURL)
	}
}

func TestFetchByTitleAuthor(t *testing.T) {
	var gotQuery string
	g := newTestGoodreads(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("q")
		w.Write([]byte(coverWithSizeSuffix))
	})

	url, err := g.FetchByTitleAuthor("Pale Blue Dot", "Carl Sagan")
	if err != nil {
		t.Fatalf("FetchByTitleAuthor() error = %v", err)
	}

	if gotQuery != "Pale Blue Dot" {
		t.Errorf("FetchByTitleAuthor() queried q = %q, want %q", gotQuery, "Pale Blue Dot")
	}

	expectedURL := "https://example.com/cover.jpg"
	if url != expectedURL {
		t.Errorf("FetchByTitleAuthor() = %v, want %v", url, expectedURL)
	}
}

func TestFetch_NotFound(t *testing.T) {
	tests := []struct {
		name    string
		fetch   func(g *Goodreads) (string, error)
		wantErr string
	}{
		{
			name:    "by ISBN",
			fetch:   func(g *Goodreads) (string, error) { return g.FetchByISBN("1234567890123") },
			wantErr: "image was not found for ISBN 1234567890123",
		},
		{
			name: "by title and author",
			fetch: func(g *Goodreads) (string, error) {
				return g.FetchByTitleAuthor("Unknown Book", "Unknown Author")
			},
			wantErr: "image was not found [book_title=Unknown+Book, author_name=Unknown+Author]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGoodreads(t, respondWith(`[]`))

			_, err := tt.fetch(g)
			if err == nil {
				t.Fatal("expected error for empty result, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestFetch_InvalidJSON(t *testing.T) {
	tests := []struct {
		name  string
		fetch func(g *Goodreads) (string, error)
	}{
		{
			name:  "by ISBN",
			fetch: func(g *Goodreads) (string, error) { return g.FetchByISBN("1234567890123") },
		},
		{
			name:  "by title and author",
			fetch: func(g *Goodreads) (string, error) { return g.FetchByTitleAuthor("Some Book", "Some Author") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGoodreads(t, respondWith(`<html>not json</html>`))

			_, err := tt.fetch(g)
			if err == nil {
				t.Fatal("expected error for invalid JSON, got nil")
			}
			if !strings.Contains(err.Error(), "failed to parse response") {
				t.Errorf("error = %q, want it to contain 'failed to parse response'", err.Error())
			}
			if strings.Contains(err.Error(), "image was not found") {
				t.Errorf("error = %q, invalid JSON must not be reported as not found", err.Error())
			}
		})
	}
}

func TestFetch_NetworkError(t *testing.T) {
	tests := []struct {
		name  string
		fetch func(g *Goodreads) (string, error)
	}{
		{
			name:  "by ISBN",
			fetch: func(g *Goodreads) (string, error) { return g.FetchByISBN("1234567890123") },
		},
		{
			name:  "by title and author",
			fetch: func(g *Goodreads) (string, error) { return g.FetchByTitleAuthor("Some Book", "Some Author") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Nothing listens on this port, so the connection is refused
			g := &Goodreads{autoCompleteURL: "http://127.0.0.1:1/?q="}

			_, err := tt.fetch(g)
			if err == nil {
				t.Fatal("expected error for refused connection, got nil")
			}
			if !strings.Contains(err.Error(), "failed to fetch URL") {
				t.Errorf("error = %q, want it to contain 'failed to fetch URL'", err.Error())
			}
		})
	}
}

func TestFetch_LargeBody(t *testing.T) {
	// Padding after the JSON array is ignored by the decoder only if the body is fully read,
	// so a large valid payload checks that the whole response is consumed.
	entries := strings.Repeat(`{"imageUrl": "https://example.com/filler.jpg"},`, 5000)
	body := `[{"imageUrl": "https://example.com/cover._SX98_.jpg"},` + entries + `{"imageUrl": "https://example.com/last.jpg"}]`
	g := newTestGoodreads(t, respondWith(body))

	url, err := g.FetchByISBN("1234567890123")
	if err != nil {
		t.Fatalf("FetchByISBN() error = %v", err)
	}

	expectedURL := "https://example.com/cover.jpg"
	if url != expectedURL {
		t.Errorf("FetchByISBN() = %v, want %v", url, expectedURL)
	}
}
