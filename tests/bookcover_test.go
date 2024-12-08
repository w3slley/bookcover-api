package tests

import (
	"bookcover-api/internal/routes"
	"testing"
)

func TestGetUrlForQuerySearch(t *testing.T) {
	data := []byte("<html></html>")
	authorName := "Author Name"
	bookTitle := "Book Title"
	key := "book+title+author+name"

	url, _ := routes.GetUrlForQuerySearch(data, bookTitle, authorName, key)

	expectedUrl := ""
	if url != expectedUrl {
		t.Errorf("GetUrlForQuerySearch() = %v, want %v", url, expectedUrl)
	}
}
