package scraper

import (
	"testing"
)

func TestNewGoodreads(t *testing.T) {
	scraper := NewGoodreads()
	if scraper == nil {
		t.Error("NewGoodreads() returned nil")
	}
}

func TestParseHTML(t *testing.T) {
	g := NewGoodreads()

	html := []byte(`
		<html>
			<body>
				<div class="test">Hello World</div>
			</body>
		</html>
	`)

	doc, err := g.parseHTML(html)
	if err != nil {
		t.Errorf("parseHTML() error = %v", err)
	}

	if doc == nil {
		t.Error("parseHTML() returned nil document")
	}

	text := doc.Find(".test").Text()
	if text != "Hello World" {
		t.Errorf("parseHTML() text = %v, want %v", text, "Hello World")
	}
}

func TestParseHTML_InvalidHTML(t *testing.T) {
	g := NewGoodreads()

	// Empty HTML should still parse successfully (goquery is forgiving)
	doc, err := g.parseHTML([]byte(""))
	if err != nil {
		t.Errorf("parseHTML() with empty HTML error = %v", err)
	}

	if doc == nil {
		t.Error("parseHTML() returned nil document for empty HTML")
	}
}

func TestExtractURLFromISBN(t *testing.T) {
	g := NewGoodreads()

	// Test with HTML that has the expected structure
	html := []byte(`
		<html>
			<body>
				<div class="BookCover__image">
					<img src="https://example.com/cover.jpg" />
				</div>
			</body>
		</html>
	`)

	url, err := g.extractURLFromISBN(html, "1234567890123")
	if err != nil {
		t.Errorf("extractURLFromISBN() error = %v", err)
	}

	expectedURL := "https://example.com/cover.jpg"
	if url != expectedURL {
		t.Errorf("extractURLFromISBN() = %v, want %v", url, expectedURL)
	}
}

func TestExtractURLFromISBN_NotFound(t *testing.T) {
	g := NewGoodreads()

	html := []byte(`<html><body><div>No book cover here</div></body></html>`)

	_, err := g.extractURLFromISBN(html, "1234567890123")
	if err == nil {
		t.Error("extractURLFromISBN() expected error for missing image, got nil")
	}

	expectedError := "image was not found for ISBN 1234567890123"
	if err.Error() != expectedError {
		t.Errorf("extractURLFromISBN() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestExtractURLFromSearch(t *testing.T) {
	g := NewGoodreads()

	html := []byte(`
		<html>
			<body>
				<table>
					<tr itemscope>
						<td>
							<img class="bookCover" src="https://example.com/cover_SX98_.jpg" />
						</td>
						<td>
							<a class="authorName">Carl Sagan</a>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`)

	url, err := g.extractURLFromSearch(html, "Pale+Blue+Dot", "Carl+Sagan")
	if err != nil {
		t.Errorf("extractURLFromSearch() error = %v", err)
	}

	// The regex _[^_]*_. removes _SX98_. (including the dot after the underscore)
	expectedURL := "https://example.com/coverjpg"
	if url != expectedURL {
		t.Errorf("extractURLFromSearch() = %v, want %v", url, expectedURL)
	}
}

func TestExtractURLFromSearch_NotFound(t *testing.T) {
	g := NewGoodreads()

	html := []byte(`<html><body><div>No results</div></body></html>`)

	_, err := g.extractURLFromSearch(html, "NonExistent+Book", "Unknown+Author")
	if err == nil {
		t.Error("extractURLFromSearch() expected error for missing book, got nil")
	}
}

func TestExtractURLFromSearch_AuthorMismatch(t *testing.T) {
	g := NewGoodreads()

	html := []byte(`
		<html>
			<body>
				<tr itemscope>
					<img class="bookCover" src="https://example.com/cover.jpg" />
					<span class="authorName">Stephen King</span>
				</tr>
			</body>
		</html>
	`)

	// Search for different author
	_, err := g.extractURLFromSearch(html, "The+Stand", "Carl+Sagan")
	if err == nil {
		t.Error("extractURLFromSearch() expected error for author mismatch, got nil")
	}
}
