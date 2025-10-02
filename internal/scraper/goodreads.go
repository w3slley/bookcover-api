package scraper

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const querySeparator = "+"

type Goodreads struct{}

func NewGoodreads() *Goodreads {
	return &Goodreads{}
}

func (g *Goodreads) FetchByTitleAuthor(bookTitle, authorName string) (string, error) {
	bookTitle = strings.ReplaceAll(bookTitle, " ", querySeparator)
	authorName = strings.ReplaceAll(authorName, " ", querySeparator)

	query := "https://www.goodreads.com/search?utf8=%E2%9C%93&q=" + bookTitle + "&search_type=books"
	body, err := g.fetchHTML(query)
	if err != nil {
		return "", err
	}

	return g.extractURLFromSearch(body, bookTitle, authorName)
}

func (g *Goodreads) FetchByISBN(isbn string) (string, error) {
	query := "https://www.goodreads.com/search?utf8=✓&query=" + isbn
	body, err := g.fetchHTML(query)
	if err != nil {
		return "", err
	}

	return g.extractURLFromISBN(body, isbn)
}

func (g *Goodreads) fetchHTML(url string) ([]byte, error) {
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

func (g *Goodreads) extractURLFromISBN(data []byte, isbn string) (string, error) {
	doc, err := g.parseHTML(data)
	if err != nil {
		return "", err
	}

	imageURL, exists := doc.Find(".BookCover__image").First().Find("img").First().Attr("src")
	if !exists {
		return "", fmt.Errorf("image was not found for ISBN %s", isbn)
	}

	return imageURL, nil
}

func (g *Goodreads) extractURLFromSearch(data []byte, bookTitle, authorName string) (string, error) {
	doc, err := g.parseHTML(data)
	if err != nil {
		return "", err
	}

	url := ""
	doc.Find("tr[itemscope]").Each(func(i int, s *goquery.Selection) {
		foundURL, urlExists := s.Find(".bookCover").First().Attr("src")

		foundAuthorName := strings.Join(strings.Fields(s.Find(".authorName").First().Text()), " ")
		foundAuthorName = strings.ReplaceAll(foundAuthorName, " ", querySeparator)

		if url == "" && urlExists && strings.EqualFold(foundAuthorName, authorName) {
			url = foundURL
		}
	})

	if url == "" {
		return "", fmt.Errorf("image was not found [book_title=%s, author_name=%s]", bookTitle, authorName)
	}

	// Remove small image indicator to retrieve bigger cover image
	imageURL := regexp.MustCompile(`_[^_]*_.`).ReplaceAllString(url, "")
	return imageURL, nil
}

func (g *Goodreads) parseHTML(data []byte) (*goquery.Document, error) {
	html := string(data)
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error creating document: %w", err)
	}
	return doc, nil
}
