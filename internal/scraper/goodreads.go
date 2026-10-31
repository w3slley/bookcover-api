package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const (
	querySeparator  = "+"
	autoCompleteURL = "https://www.goodreads.com/book/auto_complete?format=json&q="
)

var (
	errImageNotFound = errors.New("image was not found")
	imageSizeSuffix  = regexp.MustCompile(`_[^_]*_.`)
)

type Goodreads struct {
	autoCompleteURL string
}

func NewGoodreads() *Goodreads {
	return &Goodreads{autoCompleteURL: autoCompleteURL}
}

type BookJsonResponse struct {
	ImageUrl string `json:"imageUrl"`
}

func (g *Goodreads) FetchByTitleAuthor(bookTitle, authorName string) (string, error) {
	bookTitle = strings.ReplaceAll(bookTitle, " ", querySeparator)
	authorName = strings.ReplaceAll(authorName, " ", querySeparator)

	body, err := g.fetch(g.autoCompleteURL + bookTitle)
	if err != nil {
		return "", err
	}

	url, err := g.extractImageURL(body)
	if err != nil {
		return "", fmt.Errorf("%w [book_title=%s, author_name=%s]", err, bookTitle, authorName)
	}

	return url, nil
}

func (g *Goodreads) FetchByISBN(isbn string) (string, error) {
	body, err := g.fetch(g.autoCompleteURL + isbn)
	if err != nil {
		return "", err
	}

	url, err := g.extractImageURL(body)
	if err != nil {
		return "", fmt.Errorf("%w for ISBN %s", err, isbn)
	}

	return url, nil
}

func (g *Goodreads) fetch(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (g *Goodreads) extractImageURL(data []byte) (string, error) {
	var books []BookJsonResponse
	if err := json.Unmarshal(data, &books); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(books) == 0 {
		return "", errImageNotFound
	}

	return imageSizeSuffix.ReplaceAllString(books[0].ImageUrl, ""), nil
}
