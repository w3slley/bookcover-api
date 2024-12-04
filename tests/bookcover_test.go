package tests

import (
	"bookcover-api/internal/routes"
	"testing"
)

func TestGetUrlForQuerySearch(t *testing.T) {
	data := []byte("<html></html>")
	authorName := "Author Name"
	bookTitle := "Book Title"

	url, _ := routes.GetUrlForQuerySearch(data, bookTitle, authorName)

	expectedUrl := ""
	if url != expectedUrl {
		t.Errorf("GetUrlForQuerySearch() = %v, want %v", url, expectedUrl)
	}
}
