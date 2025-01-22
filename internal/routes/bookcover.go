package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"bookcover-api/internal/cache"
	"bookcover-api/internal/helpers"

	"github.com/PuerkitoBio/goquery"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	GOOGLE_BOOKS_API_KEY = "GOOGLE_BOOKS_API_KEY"
	BOOK_TITLE           = "book_title"
	AUTHOR_NAME          = "author_name"
	QUERY_SEPARATOR      = "+"
)

type HttpException struct {
	statusCode int
	message    string
}

func BuildSuccessResponse(w http.ResponseWriter, url string) []byte {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]string{"url": url})
	w.WriteHeader(200)

	return buffer.Bytes()
}

func BuildErrorResponse(w http.ResponseWriter, ex HttpException) []byte {
	data, err := json.Marshal(map[string]string{"error": ex.message})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ex)
	w.WriteHeader(ex.statusCode)
	return data
}

func BookcoverSearch(w http.ResponseWriter, r *http.Request) {
	bookTitle := r.URL.Query().Get(BOOK_TITLE)
	authorName := r.URL.Query().Get(AUTHOR_NAME)
	if bookTitle == "" || authorName == "" {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest,
			message:    helpers.MANDATORY_PARAMS_MISSING,
		}))
		return
	}
	bookTitle = strings.ReplaceAll(bookTitle, " ", QUERY_SEPARATOR)
	authorName = strings.ReplaceAll(authorName, " ", QUERY_SEPARATOR)

	cacheKey := strings.ToLower(bookTitle + QUERY_SEPARATOR + authorName)
	cachedUrl, err := cache.GetCache().Get(cacheKey)
	if err != nil {
		log.Print(err)
	}
	if cachedUrl != nil {
		log.Printf("Found cache with key %s", cacheKey)
		w.Write(BuildSuccessResponse(w, string(cachedUrl.Value)))
		return
	}

	query := "https://www.goodreads.com/search?utf8=%E2%9C%93&q=" + bookTitle + "&search_type=books"
	body, err := getBody(query)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
		}))
		return
	}

	imageUrl, err := GetUrlForQuerySearch(body, bookTitle, authorName, cacheKey)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusNotFound,
			message:    err.Error(),
		}))
		return
	}

	w.Write(BuildSuccessResponse(w, imageUrl))
}

func BookcoverByIsbn(w http.ResponseWriter, r *http.Request) {
	isbn := strings.ReplaceAll(r.PathValue("isbn"), "-", "")
	if len(isbn) != 13 {
		log.Printf("Invalid ISBN %s", isbn)
		w.Write(BuildErrorResponse(w, HttpException{statusCode: http.StatusBadRequest, message: helpers.INVALID_ISBN}))
		return
	}

	query := "https://www.goodreads.com/search?utf8=✓&query=" + isbn
	body, err := getBody(query)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
		}))
		return
	}

	imageUrl, err := GetUrlForISBNSearch(body, isbn)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusNotFound,
			message:    err.Error(),
		}))
		return
	}

	w.Write(BuildSuccessResponse(w, imageUrl))
}

func getBody(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(helpers.BOOKCOVER_NOT_FOUND)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf(helpers.ERROR_READING_BODY)
	}

	return body, nil
}

func GetUrlForISBNSearch(data []byte, isbn string) (string, error) {
	html := string(data)
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("Error creating document: %v", err)
	}
	imageUrl, exists := doc.Find(".BookCover__image").First().Find("img").First().Attr("src")
	if !exists {
		return "", fmt.Errorf("Image was not found for ISBN %s")
	}
	return imageUrl, nil
}

func GetUrlForQuerySearch(data []byte, bookTitle string, authorName string, cacheKey string) (string, error) {
	doc, err := parseHTML(data)
	if err != nil {
		return "", err
	}

	url := ""
	doc.Find("tr[itemscope]").Each(func(i int, s *goquery.Selection) {
		foundUrl, urlExists := s.Find(".bookCover").First().Attr("src")

		foundAuthorName := strings.Join(strings.Fields(s.Find(".authorName").First().Text()), " ")
		foundAuthorName = strings.ReplaceAll(foundAuthorName, " ", QUERY_SEPARATOR)
		if url == "" &&
			urlExists &&
			strings.ToLower(foundAuthorName) == strings.ToLower(authorName) {
			url = foundUrl
		}
	})
	if url == "" {
		return url, fmt.Errorf("Image was not found [book_title=%s, author_name=%s]", bookTitle, authorName)
	}
	imageUrl := regexp.MustCompile(`_[^_]*_.`).ReplaceAllString(url, "") // Remove small image indicator to retrieve bigger cover image
	if cache.GetCache() != nil {
		cache.GetCache().Set(&memcache.Item{Key: cacheKey, Value: []byte(imageUrl)})
		log.Printf("Created cache for key %s", cacheKey)

	}
	return imageUrl, nil
}

func parseHTML(data []byte) (*goquery.Document, error) {
	html := string(data)
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Error creating document: %v", err)
	}
	return doc, nil
}
